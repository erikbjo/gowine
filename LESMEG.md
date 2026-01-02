# Go(d) vin

Dette repoet inneholder tre (3) scripts laget med Go:

1. ./cmd/scrape/main.go henter viner fra Vinmonopolet sin API, og scraper [Vinmonopolet](https://www.vinmonopolet.no) og [Apertif](https://www.apertif.no) sin nettside for å se etter prisdifferenser
2. ./cmd/scrape_to_filtered/main.go filtrerer produktene etter produkter med prisdifferanse og evt. andre krav
3. ./cmd/filtered_to_readme/main.go prettyprinter den filtrerte vinen til markdown format, som vist i [README](/README.md)

Scriptene ble laget ettersom Vinmonopolet ikke offentliggjør prisendringer på andre kanaler enn presse-API'et deres, og tredjeparter som [vinpuls](https://www.vinpuls.no) kan være trege med å legge ut oversikt. Og for egen læring :)
