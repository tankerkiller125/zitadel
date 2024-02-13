// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package settings

import (
	cli_client "github.com/adlerhurst/cli-client"
	pflag "github.com/spf13/pflag"
	os "os"
)

type DomainSettingsFlag struct {
	*DomainSettings

	changed bool
	set     *pflag.FlagSet

	loginNameIncludesDomainFlag                *cli_client.BoolParser
	requireOrgDomainVerificationFlag           *cli_client.BoolParser
	smtpSenderAddressMatchesInstanceDomainFlag *cli_client.BoolParser
	resourceOwnerTypeFlag                      *cli_client.EnumParser[ResourceOwnerType]
}

func (x *DomainSettingsFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("DomainSettings", pflag.ContinueOnError)

	x.loginNameIncludesDomainFlag = cli_client.NewBoolParser(x.set, "login-name-includes-domain", "")
	x.requireOrgDomainVerificationFlag = cli_client.NewBoolParser(x.set, "require-org-domain-verification", "")
	x.smtpSenderAddressMatchesInstanceDomainFlag = cli_client.NewBoolParser(x.set, "smtp-sender-address-matches-instance-domain", "")
	x.resourceOwnerTypeFlag = cli_client.NewEnumParser[ResourceOwnerType](x.set, "resource-owner-type", "")
	parent.AddFlagSet(x.set)
}

func (x *DomainSettingsFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.loginNameIncludesDomainFlag.Changed() {
		x.changed = true
		x.DomainSettings.LoginNameIncludesDomain = *x.loginNameIncludesDomainFlag.Value
	}
	if x.requireOrgDomainVerificationFlag.Changed() {
		x.changed = true
		x.DomainSettings.RequireOrgDomainVerification = *x.requireOrgDomainVerificationFlag.Value
	}
	if x.smtpSenderAddressMatchesInstanceDomainFlag.Changed() {
		x.changed = true
		x.DomainSettings.SmtpSenderAddressMatchesInstanceDomain = *x.smtpSenderAddressMatchesInstanceDomainFlag.Value
	}
	if x.resourceOwnerTypeFlag.Changed() {
		x.changed = true
		x.DomainSettings.ResourceOwnerType = *x.resourceOwnerTypeFlag.Value
	}
}

func (x *DomainSettingsFlag) Changed() bool {
	return x.changed
}