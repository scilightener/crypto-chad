package certs

import (
	"context"
	"crypto-chad-client/internal/common"
	"crypto-chad-client/internal/domain"
	pb "crypto-chad-client/internal/grpc/certs/generated"
	"crypto-chad-lib/ec"
	"crypto-chad-lib/ecdh"
	"crypto-chad-lib/rnd"
	"crypto-chad-lib/rsa"
	"crypto-chad-lib/symencr"
	"errors"
	"log"
	"math/big"
)

const saltLength = 10

//go:generate protoc --go_out=./generated --go_opt=paths=source_relative --go-grpc_out=./generated --go-grpc_opt=paths=source_relative certs.proto
type Client struct {
	u      *domain.User
	client pb.CertsClient
}

func NewClient(user *domain.User, client pb.CertsClient) *Client {
	return &Client{
		u:      user,
		client: client,
	}
}

func (c *Client) IssueCert() {
	ecdhKeys, err := ecdh.DefaultECDH.GenerateKeys()
	if err != nil {
		panic(err)
	}
	req := c.buildIssueRequest(ecdhKeys.Public)
	signKeys, err := c.client.IssueCert(context.Background(), req)
	if err != nil {
		panic(err)
	}
	mustValidateServerSign(req.GetSalt(), signKeys.GetSign().GetCipher())

	serverPoint, err := getPointFromResponse(signKeys.GetKeys().GetServerECDH())
	if err != nil {
		panic(err)
	}
	secret := ecdh.DefaultECDH.ComputeSecret(ecdhKeys.Private, ecdh.NewPublicKey(serverPoint))

	c.u.Keys = decryptRSAKeysAES(signKeys.GetKeys(), secret)
}

func (c *Client) RetrieveCert(username string) *rsa.PublicKey {
	salt := rnd.String(saltLength)
	req := &pb.CertRetrieveReq{
		Asker:        c.u.Name,
		CiphUsername: rsa.Encrypt([]byte(username), common.ServerPubKey.E, common.ServerPubKey.N),
		Salt:         salt,
	}
	p, err := c.client.RetrieveCert(context.Background(), req)
	if err != nil {
		log.Println(err)
		return nil
	}
	mustValidateServerSign(salt, p.GetSign().GetCipher())

	e := rsa.Decrypt(p.GetPubKey().GetE(), c.u.Keys.PrivateKey.D, c.u.Keys.PrivateKey.N)
	n := rsa.Decrypt(p.GetPubKey().GetN(), c.u.Keys.PrivateKey.D, c.u.Keys.PrivateKey.N)

	key := &rsa.PublicKey{
		E: new(big.Int).SetBytes(e),
		N: new(big.Int).SetBytes(n),
	}

	return key
}

func (c *Client) buildIssueRequest(key ecdh.PublicKey) *pb.CertIssueReq {
	salt := rnd.String(saltLength)
	clientECDH := &pb.ECDHPoint{
		X: key.P.X.String(),
		Y: key.P.Y.String(),
	}
	return &pb.CertIssueReq{
		CiphUsername: rsa.Encrypt([]byte(c.u.Name), common.ServerPubKey.E, common.ServerPubKey.N),
		ClientECDH:   clientECDH,
		Salt:         salt,
	}
}

func mustValidateServerSign(expected string, got [][]byte) {
	if !common.ValidateSign(expected, got) {
		panic("server sign is not valid")
	}
}

func getPointFromResponse(point *pb.ECDHPoint) (*ec.Point, error) {
	x, ok := new(big.Int).SetString(point.X, 10)
	if !ok {
		return nil, errors.New("unable to parse ecdh point x coordinate")
	}
	y, ok := new(big.Int).SetString(point.Y, 10)
	if !ok {
		return nil, errors.New("unable to parse ecdh point y coordinate")
	}
	serverPoint := ec.NewPoint(x, y)
	if !ecdh.DefaultECDH.Curve.Contains(serverPoint) {
		return nil, errors.New("provided point is not on the expected curve")
	}
	return serverPoint, nil
}

func decryptRSAKeysAES(keys *pb.Keys, secret []byte) *rsa.Keys {
	e := symencr.Decrypt(keys.GetPub().GetE(), secret)
	n := symencr.Decrypt(keys.GetPub().GetN(), secret)
	d := symencr.Decrypt(keys.GetPriv().GetD(), secret)

	return &rsa.Keys{
		PublicKey: &rsa.PublicKey{
			E: new(big.Int).SetBytes(e),
			N: new(big.Int).SetBytes(n),
		},
		PrivateKey: &rsa.PrivateKey{
			D: new(big.Int).SetBytes(d),
			N: new(big.Int).SetBytes(n),
		},
	}
}
