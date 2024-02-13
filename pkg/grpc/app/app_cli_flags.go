// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package app

import (
	cli_client "github.com/adlerhurst/cli-client"
	pflag "github.com/spf13/pflag"
	message "github.com/zitadel/zitadel/pkg/grpc/message"
	object "github.com/zitadel/zitadel/pkg/grpc/object"
	os "os"
)

type APIConfigFlag struct {
	*APIConfig

	changed bool
	set     *pflag.FlagSet

	clientIdFlag       *cli_client.StringParser
	authMethodTypeFlag *cli_client.EnumParser[APIAuthMethodType]
}

func (x *APIConfigFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("APIConfig", pflag.ContinueOnError)

	x.clientIdFlag = cli_client.NewStringParser(x.set, "client-id", "")
	x.authMethodTypeFlag = cli_client.NewEnumParser[APIAuthMethodType](x.set, "auth-method-type", "")
	parent.AddFlagSet(x.set)
}

func (x *APIConfigFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.clientIdFlag.Changed() {
		x.changed = true
		x.APIConfig.ClientId = *x.clientIdFlag.Value
	}
	if x.authMethodTypeFlag.Changed() {
		x.changed = true
		x.APIConfig.AuthMethodType = *x.authMethodTypeFlag.Value
	}
}

func (x *APIConfigFlag) Changed() bool {
	return x.changed
}

type AppFlag struct {
	*App

	changed bool
	set     *pflag.FlagSet

	idFlag         *cli_client.StringParser
	detailsFlag    *object.ObjectDetailsFlag
	stateFlag      *cli_client.EnumParser[AppState]
	nameFlag       *cli_client.StringParser
	oidcConfigFlag *OIDCConfigFlag
	apiConfigFlag  *APIConfigFlag
	samlConfigFlag *SAMLConfigFlag
}

func (x *AppFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("App", pflag.ContinueOnError)

	x.idFlag = cli_client.NewStringParser(x.set, "id", "")
	x.stateFlag = cli_client.NewEnumParser[AppState](x.set, "state", "")
	x.nameFlag = cli_client.NewStringParser(x.set, "name", "")
	x.detailsFlag = &object.ObjectDetailsFlag{ObjectDetails: new(object.ObjectDetails)}
	x.detailsFlag.AddFlags(x.set)
	x.oidcConfigFlag = &OIDCConfigFlag{OIDCConfig: new(OIDCConfig)}
	x.oidcConfigFlag.AddFlags(x.set)
	x.apiConfigFlag = &APIConfigFlag{APIConfig: new(APIConfig)}
	x.apiConfigFlag.AddFlags(x.set)
	x.samlConfigFlag = &SAMLConfigFlag{SAMLConfig: new(SAMLConfig)}
	x.samlConfigFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *AppFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "details", "oidc-config", "api-config", "saml-config")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("details"); flagIdx != nil {
		x.detailsFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("oidc-config"); flagIdx != nil {
		x.oidcConfigFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("api-config"); flagIdx != nil {
		x.apiConfigFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if flagIdx := flagIndexes.LastByName("saml-config"); flagIdx != nil {
		x.samlConfigFlag.ParseFlags(x.set, flagIdx.Args)
	}

	if x.idFlag.Changed() {
		x.changed = true
		x.App.Id = *x.idFlag.Value
	}

	if x.detailsFlag.Changed() {
		x.changed = true
		x.App.Details = x.detailsFlag.ObjectDetails
	}

	if x.stateFlag.Changed() {
		x.changed = true
		x.App.State = *x.stateFlag.Value
	}
	if x.nameFlag.Changed() {
		x.changed = true
		x.App.Name = *x.nameFlag.Value
	}

	switch cli_client.FieldIndexes(args, "oidc-config", "api-config", "saml-config").Last().Flag {
	case "oidc-config":
		if x.oidcConfigFlag.Changed() {
			x.changed = true
			x.App.Config = &App_OidcConfig{OidcConfig: x.oidcConfigFlag.OIDCConfig}
		}
	case "api-config":
		if x.apiConfigFlag.Changed() {
			x.changed = true
			x.App.Config = &App_ApiConfig{ApiConfig: x.apiConfigFlag.APIConfig}
		}
	case "saml-config":
		if x.samlConfigFlag.Changed() {
			x.changed = true
			x.App.Config = &App_SamlConfig{SamlConfig: x.samlConfigFlag.SAMLConfig}
		}
	}
}

func (x *AppFlag) Changed() bool {
	return x.changed
}

type AppNameQueryFlag struct {
	*AppNameQuery

	changed bool
	set     *pflag.FlagSet

	nameFlag   *cli_client.StringParser
	methodFlag *cli_client.EnumParser[object.TextQueryMethod]
}

func (x *AppNameQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("AppNameQuery", pflag.ContinueOnError)

	x.nameFlag = cli_client.NewStringParser(x.set, "name", "")
	x.methodFlag = cli_client.NewEnumParser[object.TextQueryMethod](x.set, "method", "")
	parent.AddFlagSet(x.set)
}

func (x *AppNameQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if x.nameFlag.Changed() {
		x.changed = true
		x.AppNameQuery.Name = *x.nameFlag.Value
	}
	if x.methodFlag.Changed() {
		x.changed = true
		x.AppNameQuery.Method = *x.methodFlag.Value
	}
}

func (x *AppNameQueryFlag) Changed() bool {
	return x.changed
}

type AppQueryFlag struct {
	*AppQuery

	changed bool
	set     *pflag.FlagSet

	nameQueryFlag *AppNameQueryFlag
}

func (x *AppQueryFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("AppQuery", pflag.ContinueOnError)

	x.nameQueryFlag = &AppNameQueryFlag{AppNameQuery: new(AppNameQuery)}
	x.nameQueryFlag.AddFlags(x.set)
	parent.AddFlagSet(x.set)
}

func (x *AppQueryFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "name-query")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	if flagIdx := flagIndexes.LastByName("name-query"); flagIdx != nil {
		x.nameQueryFlag.ParseFlags(x.set, flagIdx.Args)
	}

	switch cli_client.FieldIndexes(args, "name-query").Last().Flag {
	case "name-query":
		if x.nameQueryFlag.Changed() {
			x.changed = true
			x.AppQuery.Query = &AppQuery_NameQuery{NameQuery: x.nameQueryFlag.AppNameQuery}
		}
	}
}

func (x *AppQueryFlag) Changed() bool {
	return x.changed
}

type OIDCConfigFlag struct {
	*OIDCConfig

	changed bool
	set     *pflag.FlagSet

	redirectUrisFlag             *cli_client.StringSliceParser
	responseTypesFlag            *cli_client.EnumSliceParser[OIDCResponseType]
	grantTypesFlag               *cli_client.EnumSliceParser[OIDCGrantType]
	appTypeFlag                  *cli_client.EnumParser[OIDCAppType]
	clientIdFlag                 *cli_client.StringParser
	authMethodTypeFlag           *cli_client.EnumParser[OIDCAuthMethodType]
	postLogoutRedirectUrisFlag   *cli_client.StringSliceParser
	versionFlag                  *cli_client.EnumParser[OIDCVersion]
	noneCompliantFlag            *cli_client.BoolParser
	complianceProblemsFlag       []*message.LocalizedMessageFlag
	devModeFlag                  *cli_client.BoolParser
	accessTokenTypeFlag          *cli_client.EnumParser[OIDCTokenType]
	accessTokenRoleAssertionFlag *cli_client.BoolParser
	idTokenRoleAssertionFlag     *cli_client.BoolParser
	idTokenUserinfoAssertionFlag *cli_client.BoolParser
	clockSkewFlag                *cli_client.DurationParser
	additionalOriginsFlag        *cli_client.StringSliceParser
	allowedOriginsFlag           *cli_client.StringSliceParser
	skipNativeAppSuccessPageFlag *cli_client.BoolParser
}

func (x *OIDCConfigFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("OIDCConfig", pflag.ContinueOnError)

	x.redirectUrisFlag = cli_client.NewStringSliceParser(x.set, "redirect-uris", "")
	x.responseTypesFlag = cli_client.NewEnumSliceParser[OIDCResponseType](x.set, "response-types", "")
	x.grantTypesFlag = cli_client.NewEnumSliceParser[OIDCGrantType](x.set, "grant-types", "")
	x.appTypeFlag = cli_client.NewEnumParser[OIDCAppType](x.set, "app-type", "")
	x.clientIdFlag = cli_client.NewStringParser(x.set, "client-id", "")
	x.authMethodTypeFlag = cli_client.NewEnumParser[OIDCAuthMethodType](x.set, "auth-method-type", "")
	x.postLogoutRedirectUrisFlag = cli_client.NewStringSliceParser(x.set, "post-logout-redirect-uris", "")
	x.versionFlag = cli_client.NewEnumParser[OIDCVersion](x.set, "version", "")
	x.noneCompliantFlag = cli_client.NewBoolParser(x.set, "none-compliant", "")
	x.complianceProblemsFlag = []*message.LocalizedMessageFlag{}
	x.devModeFlag = cli_client.NewBoolParser(x.set, "dev-mode", "")
	x.accessTokenTypeFlag = cli_client.NewEnumParser[OIDCTokenType](x.set, "access-token-type", "")
	x.accessTokenRoleAssertionFlag = cli_client.NewBoolParser(x.set, "access-token-role-assertion", "")
	x.idTokenRoleAssertionFlag = cli_client.NewBoolParser(x.set, "id-token-role-assertion", "")
	x.idTokenUserinfoAssertionFlag = cli_client.NewBoolParser(x.set, "id-token-userinfo-assertion", "")
	x.clockSkewFlag = cli_client.NewDurationParser(x.set, "clock-skew", "")
	x.additionalOriginsFlag = cli_client.NewStringSliceParser(x.set, "additional-origins", "")
	x.allowedOriginsFlag = cli_client.NewStringSliceParser(x.set, "allowed-origins", "")
	x.skipNativeAppSuccessPageFlag = cli_client.NewBoolParser(x.set, "skip-native-app-success-page", "")
	parent.AddFlagSet(x.set)
}

func (x *OIDCConfigFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args, "compliance-problems")

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	for _, flagIdx := range flagIndexes.ByName("compliance-problems") {
		x.complianceProblemsFlag = append(x.complianceProblemsFlag, &message.LocalizedMessageFlag{LocalizedMessage: new(message.LocalizedMessage)})
		x.complianceProblemsFlag[len(x.complianceProblemsFlag)-1].AddFlags(x.set)
		x.complianceProblemsFlag[len(x.complianceProblemsFlag)-1].ParseFlags(x.set, flagIdx.Args)
	}

	if x.redirectUrisFlag.Changed() {
		x.changed = true
		x.OIDCConfig.RedirectUris = *x.redirectUrisFlag.Value
	}
	if x.responseTypesFlag.Changed() {
		x.changed = true
		x.OIDCConfig.ResponseTypes = *x.responseTypesFlag.Value
	}
	if x.grantTypesFlag.Changed() {
		x.changed = true
		x.OIDCConfig.GrantTypes = *x.grantTypesFlag.Value
	}
	if x.appTypeFlag.Changed() {
		x.changed = true
		x.OIDCConfig.AppType = *x.appTypeFlag.Value
	}
	if x.clientIdFlag.Changed() {
		x.changed = true
		x.OIDCConfig.ClientId = *x.clientIdFlag.Value
	}
	if x.authMethodTypeFlag.Changed() {
		x.changed = true
		x.OIDCConfig.AuthMethodType = *x.authMethodTypeFlag.Value
	}
	if x.postLogoutRedirectUrisFlag.Changed() {
		x.changed = true
		x.OIDCConfig.PostLogoutRedirectUris = *x.postLogoutRedirectUrisFlag.Value
	}
	if x.versionFlag.Changed() {
		x.changed = true
		x.OIDCConfig.Version = *x.versionFlag.Value
	}
	if x.noneCompliantFlag.Changed() {
		x.changed = true
		x.OIDCConfig.NoneCompliant = *x.noneCompliantFlag.Value
	}
	if len(x.complianceProblemsFlag) > 0 {
		x.changed = true
		x.ComplianceProblems = make([]*message.LocalizedMessage, len(x.complianceProblemsFlag))
		for i, value := range x.complianceProblemsFlag {
			x.OIDCConfig.ComplianceProblems[i] = value.LocalizedMessage
		}
	}

	if x.devModeFlag.Changed() {
		x.changed = true
		x.OIDCConfig.DevMode = *x.devModeFlag.Value
	}
	if x.accessTokenTypeFlag.Changed() {
		x.changed = true
		x.OIDCConfig.AccessTokenType = *x.accessTokenTypeFlag.Value
	}
	if x.accessTokenRoleAssertionFlag.Changed() {
		x.changed = true
		x.OIDCConfig.AccessTokenRoleAssertion = *x.accessTokenRoleAssertionFlag.Value
	}
	if x.idTokenRoleAssertionFlag.Changed() {
		x.changed = true
		x.OIDCConfig.IdTokenRoleAssertion = *x.idTokenRoleAssertionFlag.Value
	}
	if x.idTokenUserinfoAssertionFlag.Changed() {
		x.changed = true
		x.OIDCConfig.IdTokenUserinfoAssertion = *x.idTokenUserinfoAssertionFlag.Value
	}
	if x.clockSkewFlag.Changed() {
		x.changed = true
		x.OIDCConfig.ClockSkew = x.clockSkewFlag.Value
	}
	if x.additionalOriginsFlag.Changed() {
		x.changed = true
		x.OIDCConfig.AdditionalOrigins = *x.additionalOriginsFlag.Value
	}
	if x.allowedOriginsFlag.Changed() {
		x.changed = true
		x.OIDCConfig.AllowedOrigins = *x.allowedOriginsFlag.Value
	}
	if x.skipNativeAppSuccessPageFlag.Changed() {
		x.changed = true
		x.OIDCConfig.SkipNativeAppSuccessPage = *x.skipNativeAppSuccessPageFlag.Value
	}
}

func (x *OIDCConfigFlag) Changed() bool {
	return x.changed
}

type SAMLConfigFlag struct {
	*SAMLConfig

	changed bool
	set     *pflag.FlagSet

	metadataXmlFlag *cli_client.BytesParser
	metadataUrlFlag *cli_client.StringParser
}

func (x *SAMLConfigFlag) AddFlags(parent *pflag.FlagSet) {
	x.set = pflag.NewFlagSet("SAMLConfig", pflag.ContinueOnError)

	x.metadataXmlFlag = cli_client.NewBytesParser(x.set, "metadata-xml", "")
	x.metadataUrlFlag = cli_client.NewStringParser(x.set, "metadata-url", "")
	parent.AddFlagSet(x.set)
}

func (x *SAMLConfigFlag) ParseFlags(parent *pflag.FlagSet, args []string) {
	flagIndexes := cli_client.FieldIndexes(args)

	if err := x.set.Parse(flagIndexes.Primitives().Args); err != nil {
		cli_client.Logger().Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}

	switch cli_client.FieldIndexes(args, "metadata-xml", "metadata-url").Last().Flag {
	case "metadata-xml":
		if x.metadataXmlFlag.Changed() {
			x.changed = true
			x.SAMLConfig.Metadata = &SAMLConfig_MetadataXml{MetadataXml: *x.metadataXmlFlag.Value}
		}
	case "metadata-url":
		if x.metadataUrlFlag.Changed() {
			x.changed = true
			x.SAMLConfig.Metadata = &SAMLConfig_MetadataUrl{MetadataUrl: *x.metadataUrlFlag.Value}
		}
	}
}

func (x *SAMLConfigFlag) Changed() bool {
	return x.changed
}