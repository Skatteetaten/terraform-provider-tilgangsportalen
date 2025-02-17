# Synkronisering ved endringer

Terraform provideren utvikles lokalt, i GitLab:
[Tilgangsportalen provider](https://gitlab.skead.no/datadrevet-fremtid/terraform-providers/tilgangsportalen).
Når koden er ferdig utviklet, og klar for en ny release, merges alle nødvendige
endringer inn i main, og deretter inn i release-branch (branch ved navn
release). Det er satt opp automatisk synkronisering til GitHub fra denne
branchen. Synkroniseringen er konfigurert som såkalt “push mirroring”. Det vil
si at det er GitLab som pusher/sender de nye kodeendringene til GitHub vha
GitHub sitt API. Dette gjøres i Settings->Repository->Mirroring.

Push mirroring er satt opp til kun å pushe én branch, nemlig “release”. Dette er
konfigurert med regex. For å identifisere seg mot GitHub brukes et PAT-token
laget for vår systembruker i GitHub, `skatteetaten-dataplattform-bot`. Det er
laget et klassisk PAT med scope `repo` og `workflow`. Nåværende token løper ut
**18. mai 2025**. Dersom man vil release en ny versjon etter dette, må
PAT-tokenet rulleres. Vi får en epost til felles postkasse (dataplattform(at)skatteetaten.no)
med en link for å regenerere tokenet ca en uke før utløp. Passord for innlogging
til GitHub ligger i teamets key vault.

På GitLab- siden er det konfigurert med autentisering i form av brukernavn og
passord. Brukernavnet er `skatteetaten-dataplattform-bot` og et PAT som
beskrevet over er brukt som passord. URL er
<https://github.com/Skatteetaten/terraform-provider-tilgangsportalen>. Synk
skjer automatisk ved endringer, men kan også trigges med å trykke på synk-knapp
i settings.
[Dokumentasjon av GitHub PAT tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).

Merk: Det kan ta litt tid før synken gjennomføres og `Last successful update` oppdateres.
