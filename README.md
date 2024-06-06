# Tilgangsportalen Terraform Provider

## DISCLAIMER

**Note:** This provider is based on an internal API and is of no use to others.
Please do not use this provider unless you are working within Skatteetaten. The
project is public because this is a requirement for being published to the
official Terraform Registry.

## Beskrivelse

Dette er en Terraform provider som lar deg definere roller og grupper i
Tilgangsportalen ved hjelp av Terraform (infrastruktur som kode). Provideren er
skrevet i Golang, og baserer seg på Hashicorp sitt
[Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
og Tilgangsportalen-APIet. Følgende guide fra hashicorp er brukt som
utgangspunkt og inspirasjon i utviklingen av denne provideren:
[Implement a provider with the Terraform Plugin Framework](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider).
Provideren er tilgjengelig via terraform registry:
[terraform-provider-tilgangsportalen](https://registry.terraform.io/providers/hashicorp/tilgangsportalen/latest)
Provideren er utviklet internt i Skatteetaten, av Team Dataplattform, primært
for å dekke våre egne behov. Vi håper den også vil komme andre til nytte. For
spørsmål knyttet til provideren, for å melde feil, eller for å etterspørre ny
funksjonalitet, kontakt oss via kontaktinformasjonen i kontakt-seksjonen av
denne readme-en.

## Ta i bruk provideren

Vi anbefaler at du tar i bruk provideren ved å referere til den publiserte
versjonen i Terraform Registry.

```terraform

terraform {
  required_providers {
    tilgangsportalen = {
      source  = "Skatteetaten/tilgangsportalen"
      version = "~>0.1"
    }
  }
}

provider "tilgangsportalen" {
  hosturl  = var.TILGANGSPORTALEN_URL
  username = var.TILGANGSPORTALEN_USERNAME
  password = var.TILGANGSPORTALEN_PASSWORD
}
```

## Tilgjengelig funksjonalitet

Tilgjengelig funksjonalitet er begrenset av tilgjengelige API-metoder for
Tilgangsportalen-API-et. Provideren lar deg opprette, endre, og slette
tilgangsportal-roller, opprette og slette Entra ID-grupper, og knytte disse
sammen. Det finnes også funksjonalitet for å liste ut roller, grupper, og
koblinger mellom disse. Merk at du kun vil ha mulighet til å endre eller slette
roller og grupper eid av (dette betyr som hovedregel opprettet av) den brukeren
du identifiserer deg med mot Tilgangsportalen.

## Kom i gang med utvikling

### Innstaller avhengigheter

[Go >=1.21](https://go.dev/doc/install)

`$ brew install go`

[Terraform >= 1.5+](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)

```zsh
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

For å kunne definere og endre roller og grupper i Tilgangsportalen trenger du en
gyldig bruker mot Tilgangsportalen API-et.

### Bygg av provideren

1. Klon dette repoet
2. Naviger til `terraform-provider-tilgangsportalen/internal`
3. Kjør kommandoen `go install all`

```sh
go install all

```

### Oppdatere avhengigheter

Avhengigheter kan oppdateres med:

```sh
go mod tidy

```

Bygg deretter provideren på nytt:

```sh
go install all

```

### Verifiser installasjonen

Test av provideren kjøres mot Tilgangsportalens test-api. Vi må sette lokale
variabler for URL og autentisering for å kunne kjøre lokalt: Sett hemmeligheter
(via nano ~/.zshrc):

```zsh
 export TF_VAR_TILGANGSPORTALEN_USERNAME='DittBrukernavnHer'
 export TF_VAR_TILGANGSPORTALEN_PASSWORD='DittPassordHer'
 export TF_VAR_TILGANGSPORTALEN_URL='https://tilgang-test.sits.no/ApiServer'
```

Det vil også være nødvendig å overskrive addressen for terraform provideren
under utvikling:

Opprett en .terraformrc-fil i på hjemmeområdet ditt:

```zsh
cd $HOME
nano .terraformrc
```

Legg til følgende (husk å bytte ut `<Username>`):

```terraform
provider_installation {

  dev_overrides {
      "registry.terraform.io/Skatteetaten/tilgangsportalen" = "/Users/<Username>/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Naviger deretter til
`terraform-provider-tilgangsportalen/examples/provider-install-verification`.
Endre ev. i terraform-filene der for å definere roller, assignments, outputs
e.l, og kjør deretter `terraform plan`:

```sh
terraform plan
```

Verifiser at det genereres en gyldig plan mot Tilgangsportalen.

Du kan kjøre terraform plan på samme måte for å opprette de ressursene du har
definert. Ressursene opprettes i test.

### Kjør tester

Det er lagt inn en rekke akseptansetester for provideren. Disse kjører mot
Tilgangsportalen sitt API i test. For å kjøre testene trenger du å legge inne en
test-bruker-ident. Denne må være en gyldig brukerident i Tilgangsportalen. Du
kan f.eks. sette den via nano ~/.zshrc som for de andre variablene:

```shell

 export TF_VAR_TEST_USER='a00000'
```

For at testene skal kjøre på TF_ACC være satt til en verdi. Dette kan du gjøre
når du vil kjøre tester. Naviger til mappen testene er i (de slutter på
`_test.go`) og kjør `TF_ACC=1 go test -count=1 -run='NavnPåTest' -v` for å kjøre
en spesifikk test, eller `TF_ACC=1 go test -count=1 -v` for å kjøre alle. Dersom
du vil sjekke coverage for testene, kan du kjøre `TF_ACC=1 go test -cover`.

### Debug Tilgangsportalen main client

Naviger til directory for filer (Naviger til
terraform-provider-tilgangsportalen/internal/tilgangsportalapi) og kjør :

```shell
go build .
go run .
```

## IT Shop

Opprettelse av en rolle i Tilgangsportalen er en prosess som består av to steg;

1. lag rolle
2. publiser rolle

Når man publiserer en rolle gjør man den tilgjngelig for brukere i den gitte IT
shop-en. Det finnes mulighet i Tilgangsportalen for å ha forskjellige IT shops,
slik at tilgangen for brukere til å se og søke på roller kan styres gjennom å
tilgangsstyre disse IT shop-ene. Per nå trenger vi kun å forholde oss til én IT
shop; **"General access shop shelf"** ("Access shop shelf" i test) for alle
roller.

## Kontaktinformasjon

Team Dataplattform kan nås på <dataplattform(at)skatteetaten.no>.
