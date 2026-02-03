# Synkronisering ved endringer

## Bakgrunn
Terraform provideren utvikles lokalt, i GitLab:
[Tilgangsportalen provider](https://gitlab.skead.no/datadrevet-fremtid/terraform-providers/tilgangsportalen).
Når koden er ferdig utviklet, og klar for en ny release, merges alle nødvendige
endringer inn i main, og deretter inn i release-branch (branch ved navn
release). Det er satt opp automatisk synkronisering til GitHub fra denne
branchen. Synkroniseringen er konfigurert som såkalt “push mirroring”. Det vil
si at det er GitLab som pusher/sender de nye kodeendringene til GitHub via
GitHub sitt API. Dette gjøres i Settings->Repository->Mirroring.

Push mirroring er satt opp til kun å pushe én branch, nemlig “release”. Dette er
konfigurert med regex. For å identifisere seg mot GitHub brukes et PAT-token
laget for vår systembruker i GitHub, `skatteetaten-dataplattform-bot`
Dersom man vil release en ny versjon etter dette, må PAT-tokenet rulleres.

På GitLab- siden er det konfigurert med autentisering i form av brukernavn og
passord. Brukernavnet er `skatteetaten-dataplattform-bot` og et PAT som
beskrevet over er brukt som passord.

## Rullering av Token
Token skal rulleres hver 60 dager.
Vi får en epost til felles postkasse (dataplattform(at)skatteetaten.no)
med en link for å regenerere tokenet ca en uke før utløp.
Nåværende token løper ut **09. april 2026**.

### Forutsetning
- Tilgang til GitHub Account `skatteetaten-dataplattform-bot` med passord som ligger i teamets key vault `kv-skarp-dataplat`
- Eleverte tilganger `[PAG-UC] Dataplattform - Databricks Workspace Admin - Skarp`og `[PAG-APPL] gitlab - Dataplattform - Owner` via PIM
- MFA-kode

### Instruksjon
1. Hente passord til `skatteetaten-dataplattform-bot` fra teamets key vault via Sikker VDI (Sikker Ubuntu anbefales for operasjoner i key vaults)
2. Logg inn i GitHub med `skatteetaten-dataplattform-bot` som brukernavn og passord. Spør teamlead om MFA-koden.
3. Gå til Github profilen > `Settings` > `Developer Settings` > `Personal access tokens` > `Tokens (classic)` og genere en ny klassisk token med scope for `repo`og `workflow` som utløper etter 60 dager.
[Dokumentasjon av GitHub PAT tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).
4. Ta notat av PAT, men husk å slette notat etterpå.
5. Gå til [GitLab Tilgangsportalen provider](https://gitlab.skead.no/datadrevet-fremtid/terraform-providers/tilgangsportalen) > `Settings` > `Repository` > `Mirroring repositories`.
6. Legge til nytt `Mirrored repositories`:
    1. Git repository URL: [GitHub Tilgangsportalen provider](https://github.com/Skatteetaten/terraform-provider-tilgangsportalen)
    2. Mirror direction: `Push`
    3. Authentication method: `Username and Password`
    4. Username: `skatteetaten-dataplattform-bot`
    6. Password: Sette PAT som passord
    7. Velge `Mirror specific branches` og sette den på `^release$`.
7. Husk å slette den gamle mirrored repository forbindelse.

Synk skjer automatisk ved endringer, men kan også trigges med å trykke på synk-knapp i settings.
Merk: Det kan ta litt tid før synken gjennomføres og `Last successful update` oppdateres.
