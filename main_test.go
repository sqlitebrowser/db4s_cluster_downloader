package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testCases = map[string]struct {
		url          string
		expectedType string
		expectedData string
	}{
		"currentrelease": {
			url:          "/currentrelease",
			expectedData: "3.13.1\nhttps://sqlitebrowser.org/blog/version-3-13-1-released\n",
			expectedType: "string",
		},
		"icon": {
			url:          "/favicon.ico",
			expectedData: "f546b38c57177d90c09231506100401dccf7b5b0f9f2299c3566ff132efefc96",
			expectedType: "sha256",
		},
		"indexpage": {
			url:          "/",
			expectedData: "3ea6230a3f18491d6a236aaea51863067ff0f823e6cd6eb62af5ca255ac70402",
			expectedType: "sha256",
		},

		// DB4S files

		// 3.13.1
		"DB.Browser.for.SQLite-v3.13.1-win32.msi": {
			url:          "/DB.Browser.for.SQLite-v3.13.1-win32.msi",
			expectedData: "e0b9f86d3da4d8d800e144295487e43de306c1bd27f14dccfe41e904736f25f7",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.1-win32.zip": {
			url:          "/DB.Browser.for.SQLite-v3.13.1-win32.zip",
			expectedData: "917ad2fa8d36e3bfa3fc85b11a34a8c18d189fbc2289f5a0d3bf41de8a288edc",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.1-win64.msi": {
			url:          "/DB.Browser.for.SQLite-v3.13.1-win64.msi",
			expectedData: "d023d54b3a5db10c7e896089bb3dbe6e7f4bc4eaa9bbecb34ca414be5970f688",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.1-win64.zip": {
			url:          "/DB.Browser.for.SQLite-v3.13.1-win64.zip",
			expectedData: "22375e275ec42d96de1d3b8e9ea4ed86d2a3505c4d0ffcbd1af67aa4003e5e4d",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.1-x86.64.AppImage": {
			url:          "/DB.Browser.for.SQLite-v3.13.1-x86.64-v2.AppImage",
			expectedData: "c2fd0c27c84777747527e1b28deccc824bc88eeb47f36a9575bf1ba0a5a38453",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.1.dmg": {
			url:          "/DB.Browser.for.SQLite-v3.13.1.dmg",
			expectedData: "a641cfbfcc2ce609f07de44a35134dab53485ecc18e6d9afa297b514d74bd75e",
			expectedType: "sha256",
		},

		// 3.13.0
		"DB.Browser.for.SQLite-v3.13.0-win32.msi": {
			url:          "/DB.Browser.for.SQLite-v3.13.0-win32.msi",
			expectedData: "81af0e90257b96d4ddac32b93801c160e18ec450c2727d507f80ba3c585279f3",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.0-win32.zip": {
			url:          "/DB.Browser.for.SQLite-v3.13.0-win32.zip",
			expectedData: "12c688c67acde2e76ff5d5a6c1dada854015f57c6b06c5378231fc357ddea47b",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.0-win64.msi": {
			url:          "/DB.Browser.for.SQLite-v3.13.0-win64.msi",
			expectedData: "4fd5a308481fa8ff3008bcbd069da03944698f1397b509f22a43bfda93dfccd3",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.0-win64.zip": {
			url:          "/DB.Browser.for.SQLite-v3.13.0-win64.zip",
			expectedData: "69465171d0eaca2a3d68ec5a5048f62ee192d136412a6f6747538d4535c18bbe",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.0-x86.64.AppImage": {
			url:          "/DB.Browser.for.SQLite-v3.13.0-x86.64.AppImage",
			expectedData: "58f4e35c7e8344fe1cf8f2431463b40be761c2120381257afbede2ff39fa21bc",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-v3.13.0.dmg": {
			url:          "/DB.Browser.for.SQLite-v3.13.0.dmg",
			expectedData: "dfa72811ab9faa522586a31bf680db1604442e35a2725f0aed77d5f66388724b",
			expectedType: "sha256",
		},

		"DB.Browser.for.SQLite-3.12.2-win32.msi": {
			url:          "/DB.Browser.for.SQLite-3.12.2-win32.msi",
			expectedData: "2b87a0ca1b14f436f2dc2cbfaa380249e754c3c87c81b6648a513f75d3c73368",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.2-win32.zip": {
			url:          "/DB.Browser.for.SQLite-3.12.2-win32.zip",
			expectedData: "9344bcd50865663674f11c1d8297c0d2b4a4f7ced0a459c9e71e89382549454f",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.2-win64.msi": {
			url:          "/DB.Browser.for.SQLite-3.12.2-win64.msi",
			expectedData: "723d601f125b0d2402d9ea191e4b310345ec52f76b61e117bf49004a2ff9b8ae",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.2-win64.zip": {
			url:          "/DB.Browser.for.SQLite-3.12.2-win64.zip",
			expectedData: "559edc274a2823264e886159eaa36332fd5af1f2f4b86ba2a5ef485b6420ab54",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-arm64-3.12.2.dmg": {
			url:          "/DB.Browser.for.SQLite-arm64-3.12.2.dmg",
			expectedData: "0c2076e4479cb9db5c85123cfe9750641f92566694ff9f6c99906321a2c424e8",
			expectedType: "sha256",
		},
		"DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage": {
			url:          "/DB_Browser_for_SQLite-v3.12.2-x86_64.AppImage",
			expectedData: "ea14c7439f7e666f3e9d8cbffe9048134b87db3e2d7bf65f4146b0649536de5c",
			expectedType: "sha256",
		},
		"SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe": {
			url:          "/SQLiteDatabaseBrowserPortable_3.12.2_English.paf.exe",
			expectedData: "a597b791949c260e31908d00bde474cbb4b16d55120be92ee6e0d7c08be56809",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.0-win32.msi": {
			url:          "/DB.Browser.for.SQLite-3.12.0-win32.msi",
			expectedData: "67f2bd4574fc46f0769bb6fcd940a91367cf32e56a94d4dbd6efe156dfc48e43",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.0-win32.zip": {
			url:          "/DB.Browser.for.SQLite-3.12.0-win32.zip",
			expectedData: "6a7676fb65027d7e808943d690e4211c8a0443bb32171f08827d8afae1f8d27c",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.0-win64.msi": {
			url:          "/DB.Browser.for.SQLite-3.12.0-win64.msi",
			expectedData: "0298b9e441f619f6945e8c52878171790aaefd84df349d84770cdde6a639a583",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.0-win64.zip": {
			url:          "/DB.Browser.for.SQLite-3.12.0-win64.zip",
			expectedData: "fcfba5148efe71d8717118ca56945cdeea2f55a1177553f696cbc085c934f5f3",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.12.0.dmg": {
			url:          "/DB.Browser.for.SQLite-3.12.0.dmg",
			expectedData: "4a7aaac7554c43ecec330d0631f356510dcad11e49bb01986ba683b6dfb59530",
			expectedType: "sha256",
		},
		"SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe": {
			url:          "/SQLiteDatabaseBrowserPortable_3.12.0_English.paf.exe",
			expectedData: "42e3bda299420b29bb01590d1902c7d2fd9ae89e7e446ddd12fad9c9a0446cb8",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.2-win32.msi": {
			url:          "/DB.Browser.for.SQLite-3.11.2-win32.msi",
			expectedData: "0a660c8eefdfbb8be6cf8be2abe223b0149ce8723cc1c19a36b88198be071abe",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.2-win32.zip": {
			url:          "/DB.Browser.for.SQLite-3.11.2-win32.zip",
			expectedData: "bdfcd05bf1890a3336a1091c6e9740d582167494d0010da061f9effab2243b9e",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.2-win64.msi": {
			url:          "/DB.Browser.for.SQLite-3.11.2-win64.msi",
			expectedData: "9db9d0c69c1372f09ef54599e3f87af3e28057a20c2bd6f59787d1cf16edb742",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.2-win64.zip": {
			url:          "/DB.Browser.for.SQLite-3.11.2-win64.zip",
			expectedData: "c6117e9d75bde6e0a6cbf51ee2356daa0ce41ca2dd3a6f3d1c221a36104531a0",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.2.dmg": {
			url:          "/DB.Browser.for.SQLite-3.11.2.dmg",
			expectedData: "022536d420dca87285864a4a948b699d01430721b511722bcf9c8713ab946776",
			expectedType: "sha256",
		},
		"SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe": {
			url:          "/SQLiteDatabaseBrowserPortable_3.11.2_Rev_2_English.paf.exe",
			expectedData: "552af97ee80c91b096e5268c553c8cb526022938fe550951b5ab02e45df28afc",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.1-win32.msi": {
			url:          "/DB.Browser.for.SQLite-3.11.1-win32.msi",
			expectedData: "76076d5c20240479238705f2211cad709f23c31cabe1682e2953bf6a7168b8d0",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.1-win32.zip": {
			url:          "/DB.Browser.for.SQLite-3.11.1-win32.zip",
			expectedData: "558cb41445f0bdd31605aaeb52264ae9839b9e21aa75369a51352956966700fc",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.1-win64.msi": {
			url:          "/DB.Browser.for.SQLite-3.11.1-win64.msi",
			expectedData: "ffe1f44f10d49c9d382e66b951125ae1ee10d4bce93e5a32dbb8547d6bf7122f",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.1-win64.zip": {
			url:          "/DB.Browser.for.SQLite-3.11.1-win64.zip",
			expectedData: "a648b8faffc6da3fcf761f921270de2a2871d4116e2f7baf5e3b0280a538164c",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.1v2.dmg": {
			url:          "/DB.Browser.for.SQLite-3.11.1v2.dmg",
			expectedData: "b0ee5b73b9c6305de79640f651ba59edd32c6a94c2245a2bda01ae8091a69b48",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.0-win32.msi": {
			url:          "/DB.Browser.for.SQLite-3.11.0-win32.msi",
			expectedData: "d1e28bb123ab758b476f1d1f86be5f9b0c4f4e55a72f9d6e29cfc7924adf44bb",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.0-win32.zip": {
			url:          "/DB.Browser.for.SQLite-3.11.0-win32.zip",
			expectedData: "f86a16c871394df8ae4d4f80536f2f784a3b250455642f65d352fed56384ef3a",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.0-win64.msi": {
			url:          "/DB.Browser.for.SQLite-3.11.0-win64.msi",
			expectedData: "83c8847d0f86354c53b30407fa4af96c9674711bf92c8705e2e4f33897fc9cdd",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.0-win64.zip": {
			url:          "/DB.Browser.for.SQLite-3.11.0-win64.zip",
			expectedData: "24390192ec1c48a7399d79001b69aef2f24fc8bd943128028dd0d6116e507d48",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.11.0.dmg": {
			url:          "/DB.Browser.for.SQLite-3.11.0.dmg",
			expectedData: "80d66a492ca3ed1f544d3dfea940c222059e9763280491a1d4cac8fb701e5720",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.10.1.dmg": {
			url:          "/DB.Browser.for.SQLite-3.10.1.dmg",
			expectedData: "9456e8ff081004bd16711959dcf3b5ecf9d304ebb0284c51b520d6ad1e0283ed",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.10.1-win32.exe": {
			url:          "/DB.Browser.for.SQLite-3.10.1-win32.exe",
			expectedData: "2d4ee7c846aa0c9db36cc18a5078c7c296b8eddea8f8564622fef4bc23fa4368",
			expectedType: "sha256",
		},
		"DB.Browser.for.SQLite-3.10.1-win64.exe": {
			url:          "/DB.Browser.for.SQLite-3.10.1-win64.exe",
			expectedData: "2a04eceaf32d5a96a8a7d8a91f78fdd0bc8c44a5ae7f86cde568fee27d422d12",
			expectedType: "sha256",
		},
		"SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe": {
			url:          "/SQLiteDatabaseBrowserPortable_3.10.1_English.paf.exe",
			expectedData: "bd55d13f3fd8fe82ec856cfb430e428b0d921622e0cc5ed192cb5af827bf5f77",
			expectedType: "sha256",
		},
	}
)

func TestDB4SDownloader(t *testing.T) {
	// Read the config file
	err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Don't log requests
	RecordDownloadsLocation = RECORD_NOWHERE

	// Set up Gin
	router, err := setupRouter(true)
	if err != nil {
		log.Fatal(err)
	}

	// Run through the test cases
	for name, details := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", details.url, nil)
			router.ServeHTTP(w, req)

			// Ensure the expected status code was returned
			assert.Equal(t, 200, w.Code)

			// Ensure the returned body data was correct
			switch details.expectedType {
			case "string":
				assert.Equal(t, details.expectedData, w.Body.String())
			case "sha256":
				// Calculate sha256 checksum of the body, then compare against the expected value
				s := sha256.New()
				_, err := io.Copy(s, w.Body)
				if err != nil {
					t.Errorf("Failed when calculating sha256 checksum of returned data: %s", err)
					return
				}
				shaSum := hex.EncodeToString(s.Sum(nil))
				assert.Equal(t, details.expectedData, shaSum)
			default:
				t.SkipNow()
			}
		})
	}
}
