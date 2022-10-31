# Diploma thesis

This is a repository for my diploma thesis.

## Releases

Releases are published using Github actions. Everytime there is a push to master a new release is published and contains a built pdf file. A release can be found [here](https://github.com/SamoKopecky/diploma-thesis-latex/releases). Each release is tagged with the time the push to main was made.

## Build

Install the required dependencies using [this guide](https://gist.github.com/ogajduse/ad4db70f9a6d396a133e6fd68f1a1204). After that install [latexmk](https://mg.readthedocs.io/latexmk.html).

To generate a pdf run:

```sh
latexmk -pdf thesis.pdf
```

## Info

This table displays numbers only for relevant sections (introduction, all chapters and conclusion).

| Characters | Words | Pages |
|------------|-------|-------|
| 29 834     | 4 444 | 16    |
