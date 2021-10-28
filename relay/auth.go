package relay

import (
	"context"
	"errors"
	"log"
	"time"

	cognitosrp "github.com/alexrudd/cognito-srp/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

/*
Useful links
 - https://docs.aws.amazon.com/sdk-for-go/api/service/cognitoidentityprovider/#InitiateAuthInput
 - https://github.com/alexrudd/cognito-srp
 - https://github.com/br4in3x/golang-cognito-example/blob/master/app/login.go
*/
func (r *Client) auth(ctx context.Context, username string, password string) (string, error) {
	csrp, _ := cognitosrp.NewCognitoSRP(username, password, r.poolId, r.clientId, nil)

	cfg, _ := config.LoadDefaultConfig(ctx,
		config.WithRegion(r.region),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	svc := cip.NewFromConfig(cfg)

	resp, err := svc.InitiateAuth(ctx, &cip.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		ClientId:       aws.String(csrp.GetClientId()),
		AuthParameters: csrp.GetAuthParams(),
	})
	if err != nil {
		log.Fatal(err)
	}

	if resp.ChallengeName != types.ChallengeNameTypePasswordVerifier {
		return "", errors.New("wrong challenge")
	}

	challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

	respAuth, err := svc.RespondToAuthChallenge(context.Background(), &cip.RespondToAuthChallengeInput{
		ChallengeName:      types.ChallengeNameTypePasswordVerifier,
		ChallengeResponses: challengeResponses,
		ClientId:           aws.String(csrp.GetClientId()),
	})
	if err != nil {
		return "", err
	}

	//fmt.Printf("Access Token: %s\n", *respAuth.AuthenticationResult.AccessToken)
	//fmt.Printf("ID Token: %s\n", *respAuth.AuthenticationResult.IdToken)
	//fmt.Printf("Refresh Token: %s\n", *respAuth.AuthenticationResult.RefreshToken)

	return *respAuth.AuthenticationResult.AccessToken, nil
}
