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
  url = getenv("FILE_DB_URL")
  dev = "docker://mysql/8/dev"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}