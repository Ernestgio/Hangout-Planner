data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./internal/loader/main.go",
  ]
}
env "local" {
  src = data.external_schema.gorm.url
  dev = "docker://mysql/8/dev"
  url = getenv("HANGOUT_DB_URL")
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}