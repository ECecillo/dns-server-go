# Pour tester en local

```shell
dig @127.0.0.1 -p 2053 +noedns codecrafters.io
```

Détailler la commande ?

# TODO

## Feature

### Principale

- [X] Parser la question.

- [X] Renommer les fonctions **Read** dans `header.go` et `question.go` par `Parse`.

#### Answer

- [X] Créer le fichier

- [X] Définir la structure "Answer"

- [X] Ajouter les commentaires pour documenter les champs de la structure.

- [X] Définir une fonction Read qui va lire l'objet "Answer" et print avec fmt sous forme de tableau

- [ ] Définir une fonction Write pour la structure "Answer" qui renvoie une slice de []byte.
  - [ ] On utilisera les champs qui sont dans la structure de byte qui appelle la fonction Write.

- [ ] Ajouter dans la fonction qui gère la lecture et l'écriture de la requête UDP
  - [ ] Créer un objet "Answer" avec les éléments suivants:
    - [X] Name : \x0ccodecrafters\x02io followed by a null byte (that's codecrafters.io encoded as a label sequence)
    - [X] Type : 1 encoded as a 2-byte big-endian int (corresponding to the "A" record type)
    - [X] Class : 1 encoded as a 2-byte big-endian int (corresponding to the "IN" record class)
    - [X] TTL : 	Any value, encoded as a 4-byte big-endian int. For example: 60.
    - [X] Length : 4, encoded as a 2-byte big-endian int (corresponds to the length of the RDATA field)
    - [X] Data : 	Any IP address, encoded as a 4-byte big-endian int. For example: \x08\x08\x08\x08 (that's 8.8.8.8 encoded as a 4-byte integer)

- [ ] Créer un test qui réalise ce que codecraft fait.

### Refactorisation

- [ ] Refactoriser le code en se basant sur la manière dont est implémenté le serveur DNS dans la net/dns.

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

