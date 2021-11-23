# Bird Data Estimation

Given a bird's latin and english names, scrape wikipedia
to try to ascertain the association of the bird with several
properties:

- what do they eat?
- where do they live?
- what type of nests do they build?
- what is their average wingspan?
- how many eggs do they typically lay?
- what is a fun fact about them?

This repository answers these questions by scraping wikipedia
and attempting regex searches and word associations.

Output is produced alongside an attribution that cites to the
wikipedia page that sourced the information.

Debug code has not been stripped out of this repository since
determining why a given property was guessed is a valuable
tool for debugging + explaning the results.
