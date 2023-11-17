package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func main() {
	var ip string
	var name string
	flag.StringVar(&ip, "ip", "", "IP Address, leave blank for auto")
	flag.StringVar(&name, "name", "", "Full subdomain with domain")
	flag.Parse()

	if ip == "" {
		slog.Debug("reaching out for public IP")

		resp, err := http.DefaultClient.Get("https://api.ipify.org/?format=json")
		if resp.StatusCode != 200 {
			slog.Error("unable to get ip address", "error", err)
			os.Exit(1)
		}

		defer resp.Body.Close()
		var ipResponse struct {
			IP string `json:"ip"`
		}
		if err = json.NewDecoder(resp.Body).Decode(&ipResponse); err != nil {
			slog.Error("json decoding error", "error", err)
			os.Exit(1)
		}

		ip = ipResponse.IP
		slog.Info("got IP address", "ip", ip)
	}

	slog.Debug("loading AWS config")
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("failed to load default config", "error", err)
		os.Exit(1)
	}

	slog.Debug("defining request")
	hostedZoneId := os.Getenv("AWS_HOSTED_ZONE_ID")
	if hostedZoneId == "" {
		slog.Error("failed to get AWS_HOSTED_ZONE_ID env")
		os.Exit(1)
	}

	client := route53.NewFromConfig(cfg)

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(name),
						Type: types.RRTypeA,
						TTL:  aws.Int64(600),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
					},
				},
			},
		},
		HostedZoneId: aws.String(hostedZoneId),
	}

	slog.Debug("performing request")
	if _, err = client.ChangeResourceRecordSets(ctx, input); err != nil {
		slog.Error("failed to change record set", "error", err)
		os.Exit(1)
	}
	slog.Info("success")
}
