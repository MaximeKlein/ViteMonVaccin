# Vite mon vaccin

## Configuration DB
Installer Mongodb
Rajouter un DB "ViteMonVaccin" avec deux collections VaccinationCenters, et Timeslots
Deux fichiers de backups sont fournis pour avoir quelques data dans les collections

## Lancer l'application
go build
./Vmv

## Endpoints

### Get locations
Pour voir les differents centres de vaccinations:
http://localhost:10000/get_locations

La réponse est au format JSON, et contient tous les centre de vaccinations possibles (Nom, addresse, ID). 
L'ID est réutilisé pour voir les disponibilitées dans get_availabilities 

### Get timeslots

http://localhost:10000/get_timeslots?id=0

L'id passée en parametre est celle du centre de vaccination.
La réponse est au format JSON et contient:
 - Si l'utilisateur est authentifié:
    - tous les timeslots ainsi que les gens a qui ils sont aloués
 - Sinon: 
    - tous les timeslots encore disponibles


### Book appointment

http://localhost:10000/book_apointment

Ne supporte que le POST
est attendu en parametre:
    - l'id du timeslot voulu
    - le nom de la personne qui book le timeslot
