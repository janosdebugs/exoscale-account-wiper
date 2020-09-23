
resource "random_id" "bucket" {
  byte_length = 5
}

resource "aws_s3_bucket" "test" {
  bucket = "exoscale-account-wiper-test-${replace(lower(random_id.bucket.b64_url), "_", "")}"
}

resource "aws_s3_bucket_object" "test" {
  bucket = aws_s3_bucket.test.bucket
  key = "test-${count.index}.txt"
  content = "Hello world ${count.index}!"
  count = 300
}
