# Pour tester en local

```shell
dig @127.0.0.1 -p 2053 +noedns codecrafters.io
```

Détailler la commande ?

# TODO

## Feature

### Principale

- [X] Parser la question.
- [ ] Pouvoir créer le segment pour la réponse.
- [ ] Refactoriser le code en se basant sur la manière dont est implémenté le serveur DNS dans la net/dns.

#### Réponse

- [ ] Créer le fichier

- [ ] Définir la structure "Answer"

- [ ] Ajouter les commentaires pour documenter les champs de la structure.

- [ ] Définir une fonction Write pour la structure "Answer" qui renvoie une slice de []byte.
  - [ ] On utilisera les champs qui sont dans la structure de byte qui appelle la fonction Write.
  - [ ]

- [ ] Ajouter dans la fonction qui gère la lecture et l'écriture de la requête UDP
  - [ ] Créer un objet "Answer" avec les éléments suivants:
    - [ ] Name :
    - [ ] Type :
    - [ ] Class :
    - [ ] TTL :
    - [ ] Length :
    - [ ] Data :


### Infra

- [ ] Créer un Dockerfile qui bundle le serveur.
- [ ] Pousser le Dockerfile sur le registry.
- [ ] Mettre en place le deployment de l'app sur le serveur VPS

### Additionnel

- [ ] Créer une interface pour jouer avec le serveur.

## Documentation

- [ ] Détailler en quoi consiste une serveur DNS
- [ ] Expliquer comment run le projet et tester avec `dig` par exemple.


## Test

### Server

- [ ] Créer une connexion avec le serveur UDP.
  - [ ] Envoyer un requête TCP avec un Header et s'assurer que l'ID renvoyé est le même.

### DNS Header

- [ ] Tester la construction du Header avec des valeurs incorrects d'après la spec
  - [ ] Vérfiier que l'on renvoie bien une erreur et que c'est celle que l'on attend.
- [ ] Regarder sur internet quelles sont les potentiels failles que l'on peut avoir.

