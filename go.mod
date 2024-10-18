module github.com/lemmego/lemmego

go 1.22.4

require (
	github.com/a-h/templ v0.2.771
	github.com/charmbracelet/huh v0.5.3
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/httplog/v2 v2.0.11
	github.com/lemmego/api v0.0.0-20240911190531-5368ee1ea1d4
	github.com/lemmego/migration v0.1.3
	github.com/lucsky/cuid v1.2.1
	github.com/spf13/cobra v1.8.1
	golang.org/x/tools v0.24.0
	gorm.io/gorm v1.25.11
)

replace github.com/lemmego/api => ../api

require (
	cloud.google.com/go/auth v0.6.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.8 // indirect
	cloud.google.com/go/storage v1.43.0 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/catppuccin/go v0.2.0 // indirect
	github.com/charmbracelet/bubbles v0.19.0 // indirect
	github.com/charmbracelet/bubbletea v0.27.0 // indirect
	github.com/charmbracelet/lipgloss v0.13.0 // indirect
	github.com/charmbracelet/x/ansi v0.2.2 // indirect
	github.com/charmbracelet/x/exp/strings v0.0.0-20240722160745-212f7b056ed0 // indirect
	github.com/charmbracelet/x/term v0.2.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/ggicci/owl v0.8.2 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.5 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.15.3-0.20240618155329-98d742f6907a // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.187.0 // indirect
	google.golang.org/genproto v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

require (
	cloud.google.com/go v0.115.0 // indirect
	dario.cat/mergo v1.0.0
	github.com/alexedwards/scs/v2 v2.8.0 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/ggicci/httpin v0.19.0
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/romsar/gonertia v1.3.0
	golang.org/x/crypto v0.26.0
	golang.org/x/oauth2 v0.21.0
	golang.org/x/text v0.17.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/driver/postgres v1.5.9 // indirect
)
