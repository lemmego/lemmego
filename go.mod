module github.com/lemmego/lemmego

go 1.23.3

require (
	github.com/a-h/templ v0.2.771
	github.com/lemmego/db v0.1.1 // indirect
	github.com/lemmego/migration v0.1.11
	github.com/spf13/cobra v1.8.1
)

replace (
	github.com/lemmego/api => ../api
	github.com/lemmego/auth => ../auth
	github.com/lemmego/cli => ../cli
	github.com/lemmego/db => ../db
)

require (
	cel.dev/expr v0.18.0 // indirect
	cloud.google.com/go/auth v0.10.2 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.5 // indirect
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	cloud.google.com/go/iam v1.2.2 // indirect
	cloud.google.com/go/monitoring v1.21.2 // indirect
	cloud.google.com/go/storage v1.47.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.25.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.49.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.49.0 // indirect
	github.com/alexedwards/scs/redisstore v0.0.0-20240316134038-7e11d57e8885 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/envoyproxy/go-control-plane v0.13.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/ggicci/owl v0.8.2 // indirect
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	github.com/go-chi/httplog/v2 v2.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lemmego/fsys v0.0.0-20241023132523-b7be6cd88ee9 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.32.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.57.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.57.0 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/sdk v1.32.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	google.golang.org/api v0.206.0 // indirect
	google.golang.org/genproto v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/grpc v1.68.0 // indirect
	google.golang.org/grpc/stats/opentelemetry v0.0.0-20241028142157-ada6787961b3 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
)

require github.com/lemmego/api v0.0.0-20241119171149-c5ab8bf10b81

require (
	cloud.google.com/go v0.116.0 // indirect
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/alexedwards/scs/v2 v2.8.0 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/ggicci/httpin v0.19.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f // indirect
	github.com/huandu/go-sqlbuilder v1.33.1 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/k0kubun/pp/v3 v3.4.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/romsar/gonertia v1.3.4 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
