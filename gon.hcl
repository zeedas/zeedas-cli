source = ["./build/zeedas-cli-darwin-amd64", "./build/zeedas-cli-darwin-arm64"]
bundle_id = "com.zeedas.zeedas-cli"

apple_id {
  username = "alan@zeedas.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: ZEEDAS, LLC"
}

zip {
  output_path = "zeedas-cli-darwin.zip"
}
