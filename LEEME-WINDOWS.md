# Instrucciones para Windows Server

## Archivos incluidos:
- `api-windows.exe` - El ejecutable de la API
- `.env-windows` - Configuración de la base de datos (renombrar a `.env`)
- `iniciar.bat` - Script para iniciar la API
- `docker-compose.yaml` - Configuración de base de datos
- `db/tablas.sql` - Scripts de base de datos

## Instalación rápida:

1. **Copiar archivos** a tu Windows Server (por ejemplo en `C:\ApiGo\`)

2. **Renombrar** `.env-windows` a `.env` y **configurar** tu base de datos:
   ```
   SQL_SERVER=192.168.0.2
   SQL_INSTANCE=SQLEXPRESS
   SQL_USER=consulta
   SQL_PASSWORD=consulta
   SQL_DATABASE2=REPORTES
   PORT=8080
   ```

3. **Crear las tablas** ejecutando el script `db/tablas.sql` en SQL Server

4. **Iniciar la API:**
   - Doble clic en `iniciar.bat`
   - O desde cmd: `api-windows.exe`

5. **Probar:** Abrir `http://localhost:8080` en el navegador

## Configuración:
- Puerto: Configurado en .env (PORT=8080)
- Base de datos: SQL Server según configuración en .env
- Servidor: 192.168.0.2\SQLEXPRESS
- Usuario: consulta / consulta

## Firewall:
Abrir puerto 8080 en Windows Firewall si es necesario.

¡Listo!
