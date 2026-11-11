# Tilgangsportalen API - Feilkoder per metode

Feilkodene er hentet fra [API-dokumentasjonen](https://wiki.sits.no/spaces/OIM/pages/1054582648/API+dokumentasjon) og dets undersider.

## Format på feilsvar

Vellykkede kall gir HTTP 200. Feilende kall gir alltid HTTP 500, og den faktiske feilkoden står først i `Message`:

```json
[{ "Message": "605: Object not found or no permission" }]
```

## Generelle feilkoder

| Kode | Beskrivelse                                   |
|------|-----------------------------------------------|
| 400  | Object with dependencies cannot be removed    |
| 500  | Internal Server Error                         |
| 600  | System user not authenticated                 |
| 601  | Missing mandatory property                    |
| 602  | Missing permissions to use API                |
| 603  | Field contains invalid value                  |
| 604  | Management violation                          |
| 605  | Object not found or no permission             |
| 606  | Object does not exist                         |
| 607  | Metadata violation                            |
| 608  | Occasional mandatory field                    |
| 609  | Object already exists                         |
| 610  | Object has dependencies and cannot be deleted |
| 611  | Assignment does not exist                     |
| 612  | Assignment already exists                     |

---

## Role APIs

| Metode              | Feilkoder                                   |
|---------------------|---------------------------------------------|
| CreateRole          | 500, 600, 601, 602, 603, 606, 607, 608, 609 |
| CreateSystemRole    | 500, 600, 601, 602, 603, 606, 607, 608, 609 |
| PublishRoleToITShop | 500, 602, 603, 605, 606, 608, 609           |
| RenameRole          | 500, 600, 601, 602, 605, 607, 608, 609      |
| DeleteRole          | 500, 600, 601, 602, 603, 605, 608, 610      |
| GetRole             | 500, 600, 602, 605                          |
| GetRoleID           | 500, 600, 602, 605                          |
| UpdateRole          | 400, 500, 600, 601, 602, 603, 605, 606, 608 |
| UpdateRoleL2        | 500, 600, 601, 602, 605, 606, 608           |
| UpdateRoleL3        | 500, 600, 601, 602, 605, 606, 608           |

---

## Azure AD / Entra ID Group APIs

| Metode                   | Feilkoder                              |
|--------------------------|----------------------------------------|
| CreateAzureADGroup       | 500, 600, 601, 602, 603, 606, 607, 609 |
| RenameAzureADGroup       | 500, 600, 601, 602, 605, 608, 609      |
| DeleteAzureADGroup       | 500, 600, 601, 602, 605, 608, 610      |
| GetAzureADGroupV1        | 500, 600, 602, 605, 606                |
| GetAzureADGroupID        | 500, 600, 602, 605, 606                |
| ListAzureADGroups        | 500, 600, 602                          |
| ListAzureADGroupsForRole | 500, 600, 601, 602, 605                |

---

## Active Directory Group APIs

| Metode                    | Feilkoder                              |
|---------------------------|----------------------------------------|
| CreateADGroup             | 500, 600, 601, 602, 603, 606, 607, 609 |
| AssignADEntitlementToRole | 500, 600, 602, 603, 604, 605, 608, 609 |

---

## Azure AD Entitlement APIs

| Metode                           | Feilkoder                              |
|----------------------------------|----------------------------------------|
| AssignAzureADEntitlementToRole   | 500, 600, 602, 603, 604, 605, 608, 609 |
| RemoveAzureADEntitlementFromRole | 500, 600, 601, 602, 605, 606, 608      |

---

## Role-to-Role APIs

| Metode             | Feilkoder                              |
|--------------------|----------------------------------------|
| AssignRoleToRole   | 500, 600, 601, 602, 605, 606, 609, 612 |
| RemoveRoleFromRole | 500, 600, 601, 602, 605, 611           |

---

## Terms Of Use APIs

| Metode                   | Feilkoder                         |
|--------------------------|-----------------------------------|
| GetTermsOfUseForRole     | 500, 600, 601, 602, 605, 606      |
| AssignTermsOfUseToRole   | 500, 600, 601, 602, 605, 606, 609 |
| RemoveTermsOfUseFromRole | 500, 600, 601, 602, 605, 606      |

---

## Organisasjonsroller

| Metode                                       | Feilkoder          |
|----------------------------------------------|--------------------|
| ListRolemembershipsForOrganizationIdentityV1 | 500, 601, 602, 603 |
