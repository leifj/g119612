
![coverage](https://raw.githubusercontent.com/SUNET/g119612/badges/.badges/main/coverage.svg)

## ETSI trust status lists (aka ETSI 119 612 v2)

This is a golang library implementing ETSI trust status lists. The library is meant to be used primarily to create a certificate pool for validating X509 certificates.

## Trust List in Infrastructure: General Overview:

Document for the reference:
https://github.com/EWC-consortium/eudi-wallet-rfcs/blob/main/ewc-rfc012-trust-mechanism.md#433-relying-parties

```mermaid
flowchart TD
    Issuer["Issuer"] -- Issues Credential --> Credential["Verifiable Credential"]
    Credential -- Stored in --> Wallet["Wallet Unit"]
    Wallet -- Presents Credential --> Verifier["Relying Party / Verifier"]
    Credential -- Includes --> Key["Public Key / Certificate"]
    Key -- Anchored in --> TL["Trusted List (EWC TL)"]
    Verifier -- Verifies Issuer & Credential --> TL
    Wallet -- Verifies Issuer & Credential --> TL
    Wallet -- Verifies Verifier --> TL
    Verifier -- Verifies Wallet Unit Attestation --> TL
    TL -. Must Register .-> Issuer & WalletProvider["Wallet Provider"]
    TL -. "Recommended to be Registered - Section 4.3.3" .-> Verifier

     Issuer:::actor
     Credential:::doc
     Wallet:::actor
     Verifier:::actor
     Key:::key
     TL:::trustlist
    classDef trustlist fill:#fdf6b2,stroke:#d97706,color:#92400e
    classDef actor fill:#f0f9ff,stroke:#0284c7,color:#0c4a6e
    classDef doc fill:#f3f4f6,stroke:#6b7280,color:#374151
    classDef key fill:#ecfccb,stroke:#65a30d,color:#365314
    linkStyle 1 stroke:#000000
```

## Contributing

If you want to "make gen" to re-generate the golang from the etsi XSD then you must install https://github.com/xuri/xgen first. Note that the generated code is post-processed (sed) to fix a couple of "features" in xgen that I am too lazy to pursue as bugs in xgen at this point. This stuff may change so run "make gen" at your own peril. The generated code that is known to work is commited into the repo for this reason - ymmw.

