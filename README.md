# TODO

## Feature

- [ ] Parser la question.
- [ ] Pouvoir créer le segment pour la réponse.
- [ ] Créer un Dockerfile qui bundle le serveur.
- [ ] Pousser le Dockerfile sur le registry.
- [ ] Mettre en place le deployment de l'app sur le serveur VPS
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

