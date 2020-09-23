package sos

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/janoszen/exoscale-account-wiper/plugin"
	"log"
	"sync"
)

type Plugin struct {
}

func (p *Plugin) GetKey() string {
	return "sos"
}

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p *Plugin) SetParameter(_ string, _ string) error {
	return fmt.Errorf("sos bucket deletion has no options")
}

func (p *Plugin) Run(clientFactory *plugin.ClientFactory, ctx context.Context) error {
	log.Printf("deleting SOS buckets...")

	var wg sync.WaitGroup
	poolBlocker := make(chan bool, 10)

	svc, err := clientFactory.GetS3Client("ch-gva-2")
	if err != nil {
		log.Printf("failed to establish SOS session for zone ch-gva-2 (%v)", err)
		return err
	}

	input := &s3.ListBucketsInput{}
	listBucketsOutput, err := svc.ListBucketsWithContext(ctx, input)
	if err != nil {
		log.Printf("failed to list buckets (%v)", err)

	}

	buckets := listBucketsOutput.Buckets
	for _, bucket := range buckets {
		select {
		case <-ctx.Done():
			break
		default:
		}

		wg.Add(1)
		bucketName := *bucket.Name
		go func() {
			defer wg.Done()
			poolBlocker <- true
			defer func() { <-poolBlocker }()

			log.Printf("deleting objects from bucket %s...", bucketName)
			bucketLocation, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{
				Bucket: &bucketName,
			})
			if err != nil {
				log.Printf("failed to determine bucket location for %s (%v)", bucketName, err)
				return
			}

			bucketClient, err := clientFactory.GetS3Client(*bucketLocation.LocationConstraint)
			if err != nil {
				log.Printf("failed to create client for bucket %s deletion (%v)", bucketName, err)
				return
			}

			var continuationToken *string = nil
			for {
				listObjectsInput := &s3.ListObjectsV2Input{
					Bucket:            &bucketName,
					ContinuationToken: continuationToken,
				}
				listObjectsOutput, err := bucketClient.ListObjectsV2WithContext(ctx, listObjectsInput)
				if err != nil {
					log.Printf("failed to list objects in bucket %s (%v)", bucketName, err)
					break
				}

				for _, object := range listObjectsOutput.Contents {
					objectKey := *object.Key

					log.Printf("deleting object %s from bucket %s...", objectKey, bucketName)

					_, err := bucketClient.DeleteObject(&s3.DeleteObjectInput{
						Bucket: &bucketName,
						Key:    &objectKey,
					})
					if err != nil {
						log.Printf("failed to delete object %s from bucket %s (%v).", objectKey, bucketName, err)
						return
					} else {
						log.Printf("deleted object %s from bucket %s.", objectKey, bucketName)
					}
				}

				if listObjectsOutput.NextContinuationToken == nil {
					break
				} else {
					continuationToken = listObjectsOutput.NextContinuationToken
				}
			}
			log.Printf("deleted objects from bucket %s.", bucketName)

			log.Printf("deleting multipart uploads from bucket %s...", bucketName)
			listMultiPartUploadsOutput, err := bucketClient.ListMultipartUploads(&s3.ListMultipartUploadsInput{
				Bucket: &bucketName,
			})
			if err != nil {
				log.Printf("failed to list multipart uploads in bucket %s (%v)", bucketName, err)
				return
			}
			for _, upload := range listMultiPartUploadsOutput.Uploads {
				objectKey := *upload.Key
				uploadId := *upload.UploadId

				_, err := bucketClient.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
					Bucket:   &bucketName,
					Key:      &objectKey,
					UploadId: &uploadId,
				})
				if err != nil {
					log.Printf("failed to delete multi-part upload %s for key %s in bucket %s (%v)", uploadId, objectKey, bucketName, err)
				} else {
					log.Printf("deleted multi-part upload %s for key %s in bucket %s.", uploadId, objectKey, bucketName)
				}
			}
			log.Printf("deleted multipart uploads from bucket %s.", bucketName)

			log.Printf("deleting bucket %s...", bucketName)
			_, err = bucketClient.DeleteBucket(&s3.DeleteBucketInput{
				Bucket: &bucketName,
			})
			if err != nil {
				log.Printf("failed to delete bucket %s (%v)", bucketName, err)
			} else {
				log.Printf("deleted bucket %s.", bucketName)
			}
		}()
	}

	wg.Wait()

	log.Printf("deleted SOS buckets.")
	return nil
}
