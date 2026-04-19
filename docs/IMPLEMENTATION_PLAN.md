# Wedding Invitation Platform - Implementation Plan

## 🎯 Overview

Platform undangan pernikahan digital multi-tenant dengan fitur:

- **CMS** untuk admin mengelola data undangan
- **Theme Selection** untuk memilih template undangan
- **Multi-tenant** dengan subdomain otomatis per client
- **Dynamic Invitation Sites** menggunakan htmx

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         NGINX PROXY                              │
│  (SSL Termination, Subdomain Routing, Load Balancing)           │
├─────────────────────────────────────────────────────────────────┤
│                              │                                   │
│        ┌──────────────────┬─┴─────────────────┐                 │
│        ▼                  ▼                   ▼                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐      │
│  │ CMS Frontend │  │ Client Sites │  │   Backend API    │      │
│  │   (Nuxt 3)   │  │   (htmx)     │  │    (Fiber)       │      │
│  │              │  │              │  │                  │      │
│  │ admin.domain │  │ *.domain.com │  │ api.domain.com   │      │
│  └──────────────┘  └──────────────┘  └──────────────────┘      │
│                              │                                   │
│                              ▼                                   │
│                      ┌──────────────┐                           │
│                      │  PostgreSQL  │                           │
│                      │  (Database)  │                           │
│                      └──────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
```

## 📊 Database Schema

### Core Tables

- tenant: Multi-tenant clients
- invitation: Main invitation data
- theme: Template themes
- rsvp: RSVP responses
- gallery: Image gallery
- gift_account: Digital gift/amplop
- guest_message: Guest book

## 🔄 API Endpoints

### CMS API (JSON)

- Authentication (login, logout, profile)
- Tenant Management (CRUD)
- Invitation Management (CRUD, publish/unpublish)
- Theme Management
- Gallery Management
- RSVP Management
- Guest Messages

### HTMX Endpoints (HTML Fragments)

- Main invitation page
- Couple info, Events, Gallery fragments
- RSVP form submission
- Guest messages

## 🚀 Implementation Phases

### Phase 1: Core Backend

### Phase 2: CMS Frontend

### Phase 3: Client Sites

### Phase 4: Multi-tenant & Deployment

See full documentation in each project folder.
