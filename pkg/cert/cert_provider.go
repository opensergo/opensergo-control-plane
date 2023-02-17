//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cert

import (
	"crypto/tls"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var certMU sync.RWMutex

type CertProvider interface {
	GetCert(info *tls.ClientHelloInfo) (*tls.Certificate, error)
}

// EnvCertProvider reads cert and secret from ENV
type EnvCertProvider struct {
	certEnvKey string
	pkEnvKey   string
}

func (e *EnvCertProvider) GetCert(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certMU.Lock()
	defer certMU.Unlock()
	c, s := os.Getenv(e.certEnvKey), os.Getenv(e.pkEnvKey)
	if c == "" || s == "" {
		return nil, errors.New("Read empty certificate or secret from env")
	}
	// In environment variable, the \n is replaced by whitespace character
	// So we need to replace the whitespace with \n character in the first and last line of PEM file
	// If not, the X509KeyPair func can not recognize the string from environment variable
	c, s = e.polishCert(c, s)
	keyPair, err := tls.X509KeyPair([]byte(c), []byte(s))
	return &keyPair, err
}

func (e *EnvCertProvider) polishCert(c, s string) (string, string) {
	c = strings.Replace(c, "----- ", "-----\n", -1)
	c = strings.Replace(c, " -----", "\n-----", -1)
	s = strings.Replace(s, "----- ", "-----\n", -1)
	s = strings.Replace(s, " -----", "\n-----", -1)
	return c, s
}
