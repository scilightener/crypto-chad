package certs

import (
	"context"
	"crypto-chad-lib/ec"
	"crypto-chad-lib/ecdh"
	"crypto-chad-lib/rsa"
	"crypto-chad-lib/symencr"
	pb "crypto-chad-server/internal/handlers/grpc/certs/generated"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"math/big"
)

type Handler interface {
	RetrieveCert(string) (*rsa.PublicKey, error)
	IssueCert(string) (*rsa.Keys, error)
	Keys() rsa.Keys
}

func Register(s *grpc.Server, certHandler Handler) {
	pb.RegisterCertsServer(s, &certServer{handler: certHandler})
}

//go:generate protoc --go_out=./generated --go_opt=paths=source_relative --go-grpc_out=./generated --go-grpc_opt=paths=source_relative certs.proto
type certServer struct {
	pb.UnimplementedCertsServer
	handler Handler
}

func (s *certServer) IssueCert(_ context.Context, req *pb.CertIssueReq) (*pb.SignKeys, error) {
	username := s.decryptUsername(req.GetCiphUsername())
	keys, err := s.handler.IssueCert(username)
	if err != nil {
		return nil, err
	}

	clientPoint, err := getPointFromRequest(req)
	if err != nil {
		return nil, err
	}
	ecdhKeys, err := ecdh.DefaultECDH.GenerateKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate app keys: %s", err.Error())
	}
	secret := ecdh.DefaultECDH.ComputeSecret(ecdhKeys.Private, ecdh.NewPublicKey(clientPoint))
	ciphKeys := encryptRSAKeysAES(keys, ecdhKeys.Public, secret)

	sign := s.sign(req.GetSalt())
	return &pb.SignKeys{
		Keys: ciphKeys,
		Sign: sign,
	}, nil
}

func (s *certServer) RetrieveCert(_ context.Context, req *pb.CertRetrieveReq) (*pb.SignPubKey, error) {
	username := s.decryptUsername(req.GetCiphUsername())
	pubKey, err := s.handler.RetrieveCert(username)
	if err != nil {
		return nil, err
	}
	askerPubKey, err := s.handler.RetrieveCert(req.GetAsker())
	if err != nil {
		return nil, err
	}

	ciphPubKey := cipherRSAPubKey(pubKey, askerPubKey)
	sign := s.sign(req.GetSalt())
	return &pb.SignPubKey{
		PubKey: ciphPubKey,
		Sign:   sign,
	}, nil
}

func (s *certServer) decryptUsername(ciphUsername [][]byte) string {
	return string(rsa.Decrypt(ciphUsername, s.handler.Keys().PrivateKey.D, s.handler.Keys().PrivateKey.N))
}

func (s *certServer) sign(salt string) *pb.Sign {
	key := s.handler.Keys().PrivateKey
	cipher := rsa.Encrypt([]byte(salt), key.D, key.N)
	return &pb.Sign{Cipher: cipher}
}

func getPointFromRequest(req *pb.CertIssueReq) (*ec.Point, error) {
	x, ok := new(big.Int).SetString(req.ClientECDH.X, 10)
	if !ok {
		return nil, errors.New("unable to parse ecdh point x coordinate")
	}
	y, ok := new(big.Int).SetString(req.ClientECDH.Y, 10)
	if !ok {
		return nil, errors.New("unable to parse ecdh point y coordinate")
	}
	clientPoint := ec.NewPoint(x, y)
	if !ecdh.DefaultECDH.Curve.Contains(clientPoint) {
		return nil, errors.New("provided point is not on the expected curve")
	}
	return clientPoint, nil
}

func encryptRSAKeysAES(keys *rsa.Keys, key ecdh.PublicKey, secret []byte) *pb.Keys {
	serverECDH := &pb.ECDHPoint{
		X: key.P.X.String(),
		Y: key.P.Y.String(),
	}

	ciphE := symencr.Encrypt(keys.PublicKey.E.Bytes(), secret)
	ciphN := symencr.Encrypt(keys.PublicKey.N.Bytes(), secret)
	ciphD := symencr.Encrypt(keys.PrivateKey.D.Bytes(), secret)
	return &pb.Keys{
		ServerECDH: serverECDH,
		Pub: &pb.CiphPubKey{
			E: ciphE,
			N: ciphN,
		},
		Priv: &pb.CiphPrivKey{
			D: ciphD,
			N: ciphN,
		},
	}
}

// cipherRSAPubKey encrypts keyToEncrypt rsa public key.
// encryptingKey is used to encrypt the keyToEncrypt data.
func cipherRSAPubKey(keyToEncrypt, encryptingKey *rsa.PublicKey) *pb.CiphReqPubKey {
	ciphE := rsa.Encrypt(keyToEncrypt.E.Bytes(), encryptingKey.E, encryptingKey.N)
	ciphN := rsa.Encrypt(keyToEncrypt.N.Bytes(), encryptingKey.E, encryptingKey.N)
	return &pb.CiphReqPubKey{
		E: ciphE,
		N: ciphN,
	}
}
