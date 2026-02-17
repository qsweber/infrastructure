package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ses"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		domain := "quinnweber.com"

		hostedZone, err := route53.NewZone(ctx, fmt.Sprintf("%s-zone", domain), &route53.ZoneArgs{
			Comment: pulumi.String(""),
			Name:    pulumi.String(domain),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		// SES domain identity
		domainIdentity, err := ses.NewDomainIdentity(ctx, fmt.Sprintf("%s-ses-domain-identity", domain), &ses.DomainIdentityArgs{
			Domain: pulumi.String(domain),
		})
		if err != nil {
			return err
		}

		// Verification TXT record for SES
		_, err = route53.NewRecord(ctx, fmt.Sprintf("%s-ses-domain-verification", domain), &route53.RecordArgs{
			Name:   pulumi.String(fmt.Sprintf("_amazonses.%s", domain)),
			Type:   pulumi.String("TXT"),
			ZoneId: hostedZone.ZoneId,
			Ttl:    pulumi.Int(300),
			Records: pulumi.StringArray{
				domainIdentity.VerificationToken,
			},
		})
		if err != nil {
			return err
		}

		// DKIM records
		dkim, err := ses.NewDomainDkim(ctx, fmt.Sprintf("%s-ses-domain-dkim", domain), &ses.DomainDkimArgs{
			Domain: pulumi.String(domain),
		})
		if err != nil {
			return err
		}

		for i := 0; i < 3; i++ {
			token := dkim.DkimTokens.Index(pulumi.Int(i))
			_, err := route53.NewRecord(ctx, fmt.Sprintf("%s-ses-dkim-%d", domain, i), &route53.RecordArgs{
				Name:   pulumi.Sprintf("%s._domainkey.%s", token, domain),
				Type:   pulumi.String("CNAME"),
				ZoneId: hostedZone.ZoneId,
				Ttl:    pulumi.Int(300),
				Records: pulumi.StringArray{
					pulumi.Sprintf("%s.dkim.amazonses.com", token),
				},
			})
			if err != nil {
				return err
			}
		}

		ctx.Export("zoneId", hostedZone.ZoneId)
		ctx.Export("sesDomain", pulumi.String(domain))
		ctx.Export("domainIdentityArn", domainIdentity.Arn)

		return nil
	})
}
