# 🔥 Litwick - Generador de Subtítulos Automáticos

Plataforma SaaS para generar transcripciones y subtítulos automáticos de archivos de audio/video usando AssemblyAI.

## Stack Tecnológico

### Backend
- **Go** con Fiber (framework web)
- **Supabase** (base de datos PostgreSQL + autenticación + almacenamiento)
- **GORM** (ORM)
- **AssemblyAI** (transcripción de audio)

### Frontend
- **Vue 3** con Composition API
- **Vite** (build tool)
- **Pinia** (state management)
- **Vue Router** (routing)
- **Axios** (HTTP client)

## Características - Semana 1 MVP

- ✅ Autenticación con Supabase (email/password)
- ✅ Subida de archivos de audio/video (hasta 500MB)
- ✅ Integración con AssemblyAI para transcripción
- ✅ Dashboard con lista de transcripciones
- ✅ Sistema de créditos (5 horas gratis)
- ✅ Exportar transcripciones (.txt, .srt)
- ✅ Editor de texto para corregir transcripciones
- ✅ Progress tracking del procesamiento

## Setup del Proyecto

### Pre-requisitos

1. **Go** >= 1.21
2. **Node.js** >= 18
3. Cuenta de **Supabase** (https://supabase.com) - Plan gratuito incluye BD + Storage
4. API Key de **AssemblyAI** (https://www.assemblyai.com)

### 1. Configuración de Supabase

1. **Crear proyecto**: Crea un nuevo proyecto en https://supabase.com

2. **Configurar autenticación**:
   - Ve a Authentication > Providers
   - Habilita Email Auth

3. **Crear bucket de almacenamiento**:
   - Ve a Storage
   - Crea un nuevo bucket llamado `litwick-uploads`
   - Haz el bucket público (Public bucket: ON)

4. **Copiar credenciales**:
   - **Project URL**: Settings > API > Project URL
   - **Anon Key**: Settings > API > anon/public key
   - **Service Role Key**: Settings > API > service_role key (¡secreto!)
   - **JWT Secret**: Settings > API > JWT Secret
   - **Database URL**: Settings > Database > Connection String (URI format)

### 2. Configuración de AssemblyAI

1. Regístrate en https://www.assemblyai.com
2. Obtén tu API Key (5 horas gratis al mes)
3. Copia la API Key

### 3. Backend Setup

```bash
# Instalar dependencias
go mod download

# Copiar archivo de configuración
cp .env.example .env

# Editar .env con tus credenciales
nano .env
```

**Configurar `.env`:**
```env
PORT=8080
ENVIRONMENT=development

# Usa la DATABASE_URL de Supabase (Settings > Database > Connection String)
DATABASE_URL=postgresql://postgres.xxxx:password@aws-0-us-east-1.pooler.supabase.com:5432/postgres

SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-key
SUPABASE_JWT_SECRET=your-jwt-secret

ASSEMBLYAI_API_KEY=your-assemblyai-api-key

STORAGE_BUCKET=litwick-uploads

FRONTEND_URL=http://localhost:5173
```

```bash
# Ejecutar servidor (creará las tablas automáticamente)
go run cmd/server/main.go
```

### 4. Frontend Setup

```bash
cd frontend

# Instalar dependencias
npm install

# Copiar archivo de configuración
cp .env.example .env

# Editar .env con tus credenciales
nano .env
```

**Configurar `frontend/.env`:**
```env
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
VITE_API_URL=http://localhost:8080
```

```bash
# Ejecutar servidor de desarrollo
npm run dev
```

### 5. Acceder a la Aplicación

1. Frontend: http://localhost:5173
2. Backend API: http://localhost:8080
3. Health Check: http://localhost:8080/health

## Uso

1. **Registrarse**: Crea una cuenta con email y contraseña
2. **Verificar email**: Revisa tu email para confirmar la cuenta
3. **Subir archivo**: Sube un archivo de audio/video (MP3, MP4, WAV, etc.)
4. **Procesar**: La transcripción se procesará automáticamente
5. **Ver resultados**: Accede a tu dashboard para ver las transcripciones
6. **Descargar**: Exporta en formato .txt o .srt

## Estructura del Proyecto

```
litwick/
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada del servidor
├── internal/
│   ├── config/                  # Configuración de la app
│   ├── database/                # Conexión a BD
│   ├── handlers/                # Controladores HTTP
│   ├── middleware/              # Middlewares (auth)
│   ├── models/                  # Modelos de BD
│   └── services/                # Servicios (AssemblyAI, S3, Supabase)
├── frontend/
│   ├── src/
│   │   ├── components/          # Componentes Vue
│   │   ├── views/               # Vistas/Páginas
│   │   ├── stores/              # Pinia stores
│   │   ├── router/              # Vue Router
│   │   └── config/              # Configuración frontend
│   ├── package.json
│   └── vite.config.js
├── uploads/                     # Directorio temporal
├── go.mod
├── .env.example
└── README.md
```

## API Endpoints

### Autenticación
- `GET /api/auth/me` - Obtener usuario actual

### Dashboard
- `GET /api/dashboard/` - Obtener estadísticas y transcripciones

### Transcripciones
- `GET /api/transcriptions/` - Listar transcripciones (paginado)
- `POST /api/transcriptions/upload` - Subir archivo
- `POST /api/transcriptions/:id/process` - Iniciar procesamiento
- `GET /api/transcriptions/:id` - Obtener transcripción
- `PUT /api/transcriptions/:id` - Editar texto de transcripción
- `DELETE /api/transcriptions/:id` - Eliminar transcripción
- `GET /api/transcriptions/:id/download?format=txt|srt` - Descargar

## Deploy

### Opción 1: Railway (Recomendado para monolito)

1. **Backend + BD:**
```bash
# Instalar Railway CLI
npm install -g @railway/cli

# Login
railway login

# Iniciar proyecto
railway init

# Agregar PostgreSQL
railway add -d postgresql

# Deploy
railway up
```

2. **Frontend:**
   - Conecta tu repo en Vercel
   - Configura las variables de entorno
   - Deploy automático

### Opción 2: Backend en Fly.io + Frontend en Vercel

Ver documentación específica de cada plataforma.

## Próximos Pasos (Semana 2-3)

### Semana 2
- [ ] Sistema de colas para procesamiento
- [ ] Progress bar en tiempo real (WebSockets)
- [ ] Mejorar editor de transcripciones
- [ ] Integración con Stripe/LemonSqueezy
- [ ] Landing page

### Semana 3
- [ ] Timestamps editables
- [ ] Soporte para más idiomas
- [ ] Exportar a .vtt, .docx
- [ ] Búsqueda en transcripciones
- [ ] Tests automatizados

## Créditos y Límites

- **Plan Gratuito**: 5 horas (300 minutos) al mes
- Los créditos se descuentan por minuto de audio procesado
- AssemblyAI ofrece 5 horas gratis al mes

## Troubleshooting

### Backend no inicia
- Verifica que la DATABASE_URL de Supabase sea correcta
- Verifica que todas las variables de entorno estén configuradas
- Revisa los logs en consola

### Frontend no carga
- Verifica que el backend esté corriendo en puerto 8080
- Revisa las variables de entorno en `frontend/.env`
- Revisa la consola del navegador

### Error al subir archivos
- Verifica que el bucket `litwick-uploads` exista en Supabase Storage
- Verifica que el bucket sea público
- Verifica que STORAGE_BUCKET en .env coincida con el nombre del bucket
- Revisa el tamaño del archivo (máx 500MB)

### Transcripción falla
- Verifica tu API Key de AssemblyAI
- Verifica que tengas créditos en AssemblyAI
- Revisa que el formato del archivo sea compatible

## Licencia

MIT

## Contacto

Para soporte o preguntas, abre un issue en GitHub.
