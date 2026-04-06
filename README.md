# FORUM :
## Difference entre SQL et SQLite : 
### SQL :
    Structured Query language  : est un language qui interagit vc un Systeme d gestion d bases des données (Client - Serveur ) donc MySQl est un moteur de base de données qui permet de stocker, gerer, manipuler les données ect.
    SGBD => Structure d'une base de données relationnelle (tables, relations, PK, FK)
    TYPES D COMMANDES SQL :
        DDL (data definition language) => Alter, Drop, Create  => Structure de la base 
        DML (Data Manipulation Language)  => Insert, Update, Delete => Manipulation des données 
        DQL (Data Query Language) => Select => Lecture des données 
        TCL (Transaction Control Language) => Rollback , COmmit => Gestion des Transaction 

### SQLITE: 
    est un SGBD léger intégré(EMBEDDED)
    Pas d Serveur (tout est ds un fichier local)
    une base de données = file.db => ce fichier contient des tables, données, relations 

* DATABASE.db ! 
    Creation d Fichier vide nom_BD.bd 
    Initialisation AVec HEADER "SQLite format 3" :  Version SQLite, taille des pages, infos globales.
    database.db
    ┌───────────────┐
    │ HEADER        │ (page 0)
    ├───────────────┤
    │ PAGE 1        │
    ├───────────────┤
    │ PAGE 2        │
    ├───────────────┤
    │ PAGE 3        │
    └───────────────┘

    SQLite divise le fichier en pages fixes  4096 bytes by default 
    chaque pages se sont des unités de READ / WRITE 
        TYPES des pages : 
            Table B-tree Page : Contient les lignes 
            Index B-tree Page : COntient les index
            Overflow Page : pour gros BLOBS

    - Moteur SQL recoit une requete => Tokenizer (lexer)(Découpe le texte) , Parser(analyseur comprend la structure) , Transforme en instructions
        1. Tu écris SQL
        2. SQLite (C code) reçoit string
        3. Tokenizer découpe
        4. Parser comprend
        5. Génère bytecode
        6. VM exécute
        7. B-Tree modifié
        8. Fichier .db modifié

    - SQLite = utiliser une seul structure B-tree (arbre equilibré) afin d recherche rapide + insertion rapide + données triées automatiquement 
    Une table = un B-tree ex : B-Tree (table users)
    B-Tree (logique)
    ↓
    stocké dans  Pages (physiques dans .db)

    Page 2 = racine du B-Tree users
        Si on veut inserer des données  DS CETTE TABLE, ACCEDER au page qui est dedant le B-Tree de la table users 
        Page 2 :
        [1 Sara]
        [2 Ali] RQ: si la page est pleine il va creer un nv page 
        B-Tree = plan d’organisation 
        Pages = boîtes  (c'est comme des boites  d mémoire ) : 
            ┌──────────────────────────┐
            │ 1. Header (en-tête)      │=> iNFOS DE LAPGE TYpe(table/index), Nb de lignes, Position des données 
            ├──────────────────────────┤
            │ 2. Cell pointers         │ => Pointeurs ; il trie seulement les pointeurs (chaque cell son pointeur de content du meme page  )
            ├──────────────────────────┤
            │ 3. Cell content (données)│
            └──────────────────────────┘
    
        CELL ? c'est une cellule ( ligne stockée ds SQLite), BLoc binaire representant une lgine  : =>  cell = (id=1, name='Sara')
        en interne => [rowid][taille][données encodées] 

        RQ ; le Sqlite_master = table system spécial qui une table interne ds file.db son role est decrire tous les tables, index, vues.... contient les colonnes et leur types  via colonne sql  + elle a un B-tree  c'est toujours la page 1 d file.db 

        - Difference entre LEAF et Internal : 
            LEAF : c'est la on stockent les données réelles ; => contient les cells avec lignes de la table ou les valeurs d'index+ pointeur vers ces cells ds la meme page 
            si la table a trop des données => plusieurs page etre liees en linked list 

            INTERNAL : c'est une page contient pas les données completes + contient des cells avec cle + pointeur vers une autre page(LEAF/INTERNAL) : Sert uniquement à organiser l’arbre et guider la recherche
                                 [Page 1 Internal - Root]
                                 key=5000   key=10000 ...
                                    /         |         \
                        [Page 2 Internal] [Page 3 Internal] ...
                        /       \         /       \
                    [Leaf 1]  [Leaf 2]  [Leaf 3]  [Leaf 4] ...
                    
                    RQ: ds le cas on a des millions des lignes et ona une page interne qui guider la recherche 
                        Root Internal page
                            ├─ key 5000 → Internal page 1
                            │    ├─ key 1000 → Leaf page 1 → rowid 1..1000
                            │    ├─ key 2000 → Leaf page 2 → rowid 1001..2000
                            │    ...
                            ├─ key 10000 → Internal page 2
                            │    ├─ Leaf page 101..200 → rowid 5001..10000
                            │    ...
        -Overflow page = si une ligne est trés grosse (BLob)

        -Index B-Tree : 
            un index = structure qui permet de trouver rapidement des lignes ds une table 
            l'index creer un B-tree separe 
            [key de l’index] + [rowid correspondant à la ligne de la table]

            B-Tree index stocke les valeurs reelles de la colonne comme KEY 
                key = bob@example.com → valeur du champ email
                rowid = 2 → pointeur vers la ligne complète dans la table
                EX:

                    SELECT * FROM users WHERE email='bob@example.com';
                        SQLite regarde l’index idx_email
                        Cherche dans le B-Tree la key = 'bob@example.com'
                        Trouve rowid = 2 → va lire la ligne complète dans la table
                        SQLite B-Tree Structure

[Table users B-Tree]
Internal Page (si table > 1 page)
 └─ guides rowid
Leaf Pages:
 ├─ rowid 1 → Alice
 ├─ rowid 2 → Bob
 └─ rowid 3 → Carol

[Index idx_email B-Tree]
Internal Page (si index > 1 page)
 └─ guides keys
Leaf Pages:
 ├─ key='alice@example.com' → rowid=1
 ├─ key='bob@example.com'   → rowid=2
 └─ key='carol@example.com' → rowid=3

Requête: email='bob@example.com'
  ↓
Index B-Tree → rowid=2 → Table B-Tree → ligne complète


   



