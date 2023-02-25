package router

import (
	bolt "go.etcd.io/bbolt"
	"crypto/x509"
	"reflect"
)

const MAP_BACKENDS string = "BACKENDS"

type Router struct{
	db *bolt.DB
}

func New(pathToDB string) *Router {
	db, err := bolt.Open(pathToDB, 0666, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(MAP_BACKENDS))
		if err != nil {
			panic(err)
		}
		return nil
	})

	return &Router{db}
}

func (r *Router) Close() {
	r.db.Close()
}

func (r *Router) GenerateBackendCode(certificate *x509.Certificate) string {
	var backendCode string

	r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MAP_BACKENDS))

		for {
			backendCode = generateRandomCode()
			value := bucket.Get([]byte(backendCode))
			if value == nil {
				err := bucket.Put([]byte(backendCode), certificate.Raw)
				return err
			}
		}
	})

	return backendCode
}

func (r *Router) VerifyBackendCode(requestCertificate *x509.Certificate, backendCode string) bool {
	var knownCertificate []byte
	r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(MAP_BACKENDS))
		knownCertificate = bucket.Get([]byte(backendCode))
		return nil
	})

	if reflect.DeepEqual(knownCertificate, requestCertificate.Raw) {
		return true
	}

	return false
}
