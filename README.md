 # Pipit Client
 
 Cliente web para el sistema de mensajería **Pipit**.
 
 ## ¿Qué es?
 
 Una página de chat muy simple que se conecta a la API de Pipit a través del balanceador, junto con el pequeño servidor en Go que la entrega.
 
 * El cliente **no tiene lógica de servidor propia**: todo el trabajo (guardar mensajes, consultar el historial, entregar mensajes nuevos) lo hace la API. La página solo muestra la conversación y envía lo que el usuario escribe.
 
 ```
 Navegador ──► Pipit Client ──► Pipit Balancer ──► Pipit API ──► Pipit Shards
 ```
 
 ## Componentes principales
 
 * **`main.go`:** sirve la página y hace de *proxy* hacia el balanceador. Todo lo que llega a `/api/...` se reenvía a la API, de modo que el navegador habla siempre con un único origen y no hace falta configurar CORS. Los streams de eventos pasan sin buffering.
 * **`index.html`:** el chat. Se embebe en el binario con `go:embed`, por lo que no necesita archivos adicionales.
 
 ## Cómo funciona la página
 
 1. El usuario elige un **nombre** y un **tópico** (por defecto, `general`) y presiona *Entrar*.
 2. La página abre una conexión de *Server-Sent Events* (`/topics/{topic}/stream`) para recibir los mensajes nuevos en tiempo real.
 3. Una vez conectada, carga los últimos 50 mensajes del tópico. Se abre primero el stream y después se pide el historial, así no se pierde ningún mensaje entre ambos; los duplicados se descartan por ID.
 4. Al enviar un mensaje, hace un `POST` a la API. El mensaje aparece cuando vuelve por el stream, lo que puede tardar hasta 1 segundo.
 5. Si la conexión se corta, el navegador reconecta solo y retoma desde el último mensaje recibido (`Last-Event-ID`). El estado de la conexión se muestra bajo el título.
 
 ## Uso
 
 Requiere **Go 1.22 o superior**.
 
 ```
 go run . -listen :3000 -api http://localhost:8080
 ```
 
 | Flag | Por defecto | Descripción |
 |---|---|---|
 | `-listen` | `:3000` | Dirección donde escucha el cliente. |
 | `-api` | `http://localhost:8080` | URL del balanceador de Pipit. |
 
 Luego abre `http://localhost:3000`. Para ver la conversación entre dos usuarios, abre una segunda pestaña con otro nombre y el mismo tópico.
 
 El sistema completo debe estar corriendo antes de usar el chat: primero el servidor de shards, luego las instancias de la API y por último el balanceador.
 
 ## Limitaciones
 
 * **Sin autenticación:** el nombre no se verifica; cualquiera puede escribir como otro usuario.
 * **Retardo de hasta 1 segundo** entre que se envía un mensaje y aparece en pantalla, por el intervalo con que la API consulta las shards.
 * **Fuente externa:** la tipografía (Inter) se carga desde Google Fonts. Sin conexión a internet, el navegador usa la fuente del sistema.
 
 ## Contexto
 
 Software desarrollado como proyecto de ejemplo para la asignatura **Desarrollo de Aplicaciones para Ambientes Distribuidos** del **IES6017** (Salta, Argentina).
 
 > Este software se distribuye sin garantías de ningún tipo.
