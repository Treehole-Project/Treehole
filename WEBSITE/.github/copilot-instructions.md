# Copilot Instructions

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

## Project Context

This is a Vue.js tree hole website project with the following characteristics:

- Vue 3 with Composition API
- Vue Router for navigation
- Axios for API communication
- Tailwind CSS for styling
- Modern, clean, and minimalist design

## Backend API

The backend provides REST APIs at `http://localhost:8080/api/v1/` with the following endpoints:

- GET /posts - Get posts list with pagination
- GET /posts/:id - Get single post details
- GET /posts/:id/replies - Get post replies
- GET /search - Basic search posts
- GET /search/advanced - Advanced search with multiple criteria
- GET /tags - Get available tags
- POST /posts - Create new post
- POST /posts/:id/replies - Create reply

## Code Style Guidelines

- Use Composition API with `<script setup>`
- Use TypeScript-style prop definitions where applicable
- Keep components small and focused
- Use Tailwind CSS classes for styling
- Follow Vue 3 best practices
- Use proper error handling for API calls
- Implement responsive design
- Use semantic HTML elements
