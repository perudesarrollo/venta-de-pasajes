###Servicio de gestión de equipaje

1. **Creación del paquete**: En el directorio del proyecto, crea un nuevo directorio llamado `baggage` para este microservicio.

2. **Definición de los modelos de datos**: Dentro del paquete `baggage`, crea un archivo llamado `model.go`. Aquí definiremos los modelos de datos necesarios para gestionar el equipaje. Por ejemplo, podríamos tener un modelo de datos para la reserva de equipaje y otro para los tipos de equipaje disponibles.

3. **Lógica de negocio**: Implementaremos la lógica necesaria para agregar equipaje a una reserva existente. Esto implicará escribir funciones que manejen la lógica de negocio específica para el manejo del equipaje, como el cálculo de precios adicionales, la validación de la cantidad de equipaje, etc.

4. **Integración con MongoDB**: Al igual que en el primer microservicio, necesitaremos interactuar con MongoDB para almacenar y recuperar datos relacionados con el equipaje. Implementaremos funciones para interactuar con la base de datos y almacenar la información relacionada con el equipaje.

5. **API Endpoints**: Expondremos endpoints de API para permitir que los clientes del servicio interactúen con él. Estos endpoints permitirán a los usuarios agregar equipaje a sus reservas existentes.

###Métodos

1. **CreateReservation**: crea una nueva reserva de equipaje en la base de datos.

2. **GetBaggageTypesByName**: Este método busca tipos de equipaje por nombre. Si el nombre está vacío, devuelve todos los tipos de equipaje.

3. **AddBaggageToReservation**: Este método agrega equipaje a una reserva existente. Ya estamos obteniendo el precio total del equipaje en la función `calculateBaggagePrice`.

4. **CalculateBaggagePrice**: Este método calcula el precio total del equipaje en función del tipo y la cantidad.

