package sg

import (
	"context"
	"github.com/exoscale/egoscale"
	"log"
	"sync"
)

func (p *Plugin) GetKey() string {
	return "sg"
}

func (p *Plugin) Run(client *egoscale.Client, ctx context.Context) error {
	log.Printf("deleting security groups...")

	var wg sync.WaitGroup
	poolBlocker := make(chan bool, 10)

	sg := &egoscale.SecurityGroup{}
	sgs, err := client.ListWithContext(ctx, sg)
	if err != nil {
		return err
	}
	for _, sg := range sgs {
		securityGroup := sg.(*egoscale.SecurityGroup)
		ingressRules := securityGroup.IngressRule
		egressRules := securityGroup.EgressRule
		wg.Add(1)
		go func() {
			defer wg.Done()
			poolBlocker <- true
			defer func() { <-poolBlocker }()

			log.Printf("removing rules from security group %s...", securityGroup.Name)
			for _, ingressRule := range ingressRules {
				err := client.BooleanRequestWithContext(ctx, egoscale.RevokeSecurityGroupIngress{
					ID: ingressRule.RuleID,
				})
				if err != nil {
					log.Printf(
						"failed to remove ingress rule %s from security group %s (%v)",
						ingressRule.RuleID,
						ingressRule.SecurityGroupName,
						err,
					)
					continue
				}
			}

			for _, egressRule := range egressRules {
				err := client.BooleanRequestWithContext(ctx, egoscale.RevokeSecurityGroupEgress{
					ID: egressRule.RuleID,
				})
				if err != nil {
					log.Printf(
						"failed to remove ingress rule %s from security group %s (%v)",
						egressRule.RuleID,
						egressRule.SecurityGroupName,
						err,
					)
					continue
				}
			}
			log.Printf("removed rules from security group %s.", securityGroup.Name)
		}()
	}

	wg.Wait()

	for _, sg := range sgs {
		securityGroup := sg.(*egoscale.SecurityGroup)
		sgId := securityGroup.ID
		if securityGroup.Name == "default" {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			poolBlocker <- true
			defer func() { <-poolBlocker }()

			log.Printf("deleting security group %s...", securityGroup.Name)
			err := client.BooleanRequestWithContext(ctx, egoscale.DeleteSecurityGroup{
				ID: sgId,
			})
			if err != nil {
				log.Printf("failed to delete security group %s (%v)", securityGroup.Name, err)
				return
			}
			log.Printf("deleted security group %s.", securityGroup.Name)
		}()
	}
	wg.Wait()

	log.Printf("deleted security groups.")
	return nil
}
