### Diseño de la arquitectura
1. Identificar los microservicios
Para este proyecto de venta de pasajes, identificaremos los siguientes microservicios:
* Servicio de búsqueda y reserva de ruta: 
    Este microservicio se encargará de buscar y reservar las rutas disponibles según los criterios proporcionados por el usuario.
* Servicio de gestión de equipaje: 
    Manejará la adición de equipaje al proceso de reserva de pasajes.
* Servicio de selección de asientos: 
    Permitirá al usuario elegir los asientos deseados para su viaje.
* Servicio de pago: 
    Gestionará el proceso de pago de la reserva de pasajes.
* Servicio de emisión de boletos: 
    Emitirá el boleto una vez que la compra haya sido confirmada

2. Desarrollo microservicios
*    Servicio de búsqueda y reserva de ruta
*    Servicio de búsqueda y reserva de ruta

```
venta-de-pasajes/
├── config/
│   └── config.go
├── cmd/
│   ├── search-service/
│   │   ├── Dockerfile
│   │   └── main.go
│   └── baggage-service/
│       ├── Dockerfile
│       └── main.go
├── haproxy/
│   └── haproxy.cfg
├── internal/
│   ├── search/
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── repository/
│   │       ├── mongodb_repository.go
│   │       └── mysql_repository.go
│   └── baggage/
│       ├── handler.go
│       ├── model.go
│       └── repository.go
├── docker-compose.yml
├── kubectl-archivo.yaml
├── go.sum
├── go.sum
└── README.md
```

La arquitectura basada en microservicios de `venta-de-pasajes` se compone de varios componentes distribuidos en diferentes directorios:

1. **Configuración (`config/`):** Contiene el archivo `config.go`, que probablemente define la configuración global de la aplicación.

2. **Servicios (`cmd/`):** Incluye dos servicios principales, `search-service` y `baggage-service`, cada uno con su propio directorio:
   - `search-service`: Contiene el `Dockerfile` y el archivo `main.go` para el servicio de búsqueda de pasajes.
   - `baggage-service`: Similar a `search-service`, contiene el `Dockerfile` y el archivo `main.go` para el servicio de manejo de equipaje.

3. **HAProxy (`haproxy/`):** Contiene el archivo de configuración `haproxy.cfg` para el balanceador de carga HAProxy, que probablemente enruta el tráfico entre los servicios `search-service` y `baggage-service`.

4. **Componentes internos (`internal/`):** Contiene la lógica interna de los servicios, separada por funcionalidad:
   - `search/`: Contiene los archivos relacionados con el servicio de búsqueda, incluyendo manejadores, modelos y repositorios.
   - `baggage/`: Similar a `search/`, contiene los archivos relacionados con el servicio de manejo de equipaje.

5. **Docker Compose (`docker-compose.yml`):** Define la configuración para orquestar y ejecutar los servicios de la aplicación utilizando Docker Compose. Incluye la configuración de los servicios, las redes y los volúmenes necesarios para su funcionamiento.

6. **Kubernetes (`kubectl-archivo.yaml`):** Contiene la configuración para desplegar la aplicación en un clúster de Kubernetes.

Esta arquitectura sigue el enfoque de microservicios, donde cada servicio se encarga de una funcionalidad específica y se comunica con otros servicios a través de API REST u otros mecanismos de comunicación. La separación en diferentes directorios y componentes facilita el desarrollo, la implementación y el mantenimiento de la aplicación, permitiendo escalabilidad y flexibilidad.


### docker-compose.yml
```bash
# Para construir y ejecutar la aplicación usando Docker Compose
docker-compose up --build
# seed routes
go run scripts/seedRoutes.go
```

- **haproxy:**  
  - Imagen: haproxy:latest
  - Volumen: Monta el archivo de configuración `haproxy.cfg` en el contenedor.
  - Puertos: Mapea el puerto 3000 del host al puerto 80 del contenedor.
  - Redes: Se conecta a la red `network-venta-de-pasajes`.
  - Dependencias: Dependiente de los servicios `search-service` y `baggage-service`.

- **search-service:**  
  - Build: Construye la imagen del servicio a partir del Dockerfile ubicado en `cmd/search-service/Dockerfile`.
  - Puertos: Expone el puerto 8080 del contenedor.
  - Variables de entorno: Define la variable de entorno `MONGO_URL` para especificar la URL de la base de datos MongoDB.
  - Redes: Se conecta a la red `network-venta-de-pasajes`.

- **baggage-service:**  
  - Build: Construye la imagen del servicio a partir del Dockerfile ubicado en `cmd/baggage-service/Dockerfile`.
  - Puertos: Expone el puerto 8081 del contenedor.
  - Variables de entorno: Define la variable de entorno `MONGO_URL` para especificar la URL de la base de datos MongoDB.
  - Redes: Se conecta a la red `network-venta-de-pasajes`.

- **mongodb:**  
  - Imagen: Utiliza la imagen de MongoDB.
  - Volumen: Monta un volumen para almacenar los datos de MongoDB en `~/Sites/Data/mongodb-data`.
  - Puertos: Mapea el puerto 27017 del host al puerto 27017 del contenedor.
  - Redes: Se conecta a la red `network-venta-de-pasajes`.

- **networks:**  
  - network-venta-de-pasajes: Define una red de tipo puente para conectar los servicios entre sí.

Este archivo permite definir la estructura y las relaciones entre los distintos componentes de la aplicación, facilitando su despliegue y ejecución mediante un solo comando: `docker-compose up --build`.


