# Realease og publisering av provider

## Release

Oppdater [CHANGELOG.md](../changelog.md) med ny versjon og beskrivelse av endringene du har gjort.

Lag MR til main branch, men **ikke** huk av *Squash commits when merge request is accepted*. Om man huker av for dette valget, så vil man få merge conflict ved neste merge til main branch.

Når endringene er merget til main lager du en MR mot release branch, her kan du huke av på å squashe commits.

Når endringene merges til release branch vil det trigge en [GitHub Actions workflow](https://github.com/Skatteetaten/terraform-provider-tilgangsportalen/blob/release/.github/workflows/release.yml) som synkroniserer endringene til GitHub, mer informasjon om hvordan denne prosessen fungerer kan du lese om [her](https://github.com/Skatteetaten/terraform-provider-tilgangsportalen/blob/release/doc/synkronisering-til-github.md).
Når endringene er synkronisert til GitHub, må det opprettes en ny tag.

Tagging kan gjøres i git med:

```shell
git tag v0.2.1
```

Hvor v0.2.1 er byttet ut med ønsket versjonsnummer. Det er viktig at semantic
versioning følges (<https://semver.org>).

Alternativt kan en ny tag opprettes via GUIet. Naviger til repoet og velg tags,
rett under repo-navnet. Bytt fana til Release, og trykk på “Draft a new
release”.

Når en ny tag opprettes vil GitHub workflow kjøres automatisk. Under
Actions-fanen i repoet i GitHub vil du kunne sjekke om workflow-en har kjørt OK.
Den nye releasen vil være tilgjengelig på repoets oversiktsside.

For utfyllende dokumentasjon om release i GitHub, se
<https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository>.

## Publisering

Nye versjoner vil bli publisert til Skatteetaten sitt namespace i Terraform
Register automatisk. Til dette brukes GPG secret key som ligger i repoet vårt
som en GitHub Actions secret sammen med valg passphrase for nøkkelen, som også
ligger i GitHub someen secret. For opprettelse av GPG nøkler er
[hashicorp sin guide benyttet](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-release-publish).
Type er `RSA and RSA`, key size er 4096, og det er ikke satt en utløpsdato.
Brukeridenten brukt er <dataplattform(at)skatteetaten.no>.
