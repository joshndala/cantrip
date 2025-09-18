# CanTrip Frontend

A beautiful Next.js chatbot frontend for the CanTrip AI Canadian Travel Assistant.

## Features

- **Dynamic City-Based Theming**: Automatically detects cities from chat messages and applies unique color palettes
- **Smooth Animations**: Framer Motion powered transitions and micro-interactions
- **Modern UI**: Built with Tailwind CSS and shadcn/ui components
- **Responsive Design**: Works seamlessly across all device sizes
- **Real-time Chat**: Connects to the CanTrip backend API for intelligent travel assistance

## City Themes

The application includes custom color palettes for major Canadian cities:

- **Toronto**: Urban Blue + Lake Blue + Concrete Grey
- **Vancouver**: Forest Green + Lake Blue + Sunrise Gold  
- **Montreal**: Maple Red + Sunrise Gold + Forest Green
- **Calgary**: Sunrise Gold + Maple Red + Grey
- **Ottawa**: Urban Blue + Maple Red + Concrete Grey
- **Edmonton**: Forest Green + Sunrise Gold + Grey
- **Halifax**: Lake Blue + Urban Blue + Sunrise Gold
- **Banff**: Forest Green + Lake Blue + Snow White
- **Whistler**: Urban Blue + Lake Blue + Grey
- **Jasper**: Forest Green + Lake Blue + Neutral Grey
- And many more...

## Tech Stack

- **Next.js 15** with App Router
- **TypeScript** for type safety
- **Tailwind CSS** for styling
- **shadcn/ui** for UI components
- **Framer Motion** for animations
- **Lucide React** for icons

## Getting Started

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser

## Backend Integration

The frontend connects to the CanTrip backend API running on `http://localhost:8080`. Make sure the backend is running before testing the chat functionality.

## Project Structure

```
src/
├── app/
│   ├── globals.css          # Global styles with custom CanTrip colors
│   ├── layout.tsx           # Root layout component
│   └── page.tsx             # Main chatbot page
├── components/
│   ├── ui/                  # shadcn/ui components
│   └── Chat.tsx             # Main chat component
└── lib/
    ├── cityThemes.ts        # City color palette definitions
    └── utils.ts             # Utility functions
```

## Customization

### Adding New Cities

To add a new city theme, update the `cityThemes` object in `src/lib/cityThemes.ts`:

```typescript
export const cityThemes: Record<string, CityTheme> = {
  // ... existing themes
  'new-city': {
    primary: '#3A5A8C',
    secondary: '#1E88E5', 
    accent: '#FFC107',
    background: 'linear-gradient(135deg, #3A5A8C 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-urban-blue via-lake-blue to-sunrise-gold'
  }
};
```

### Modifying Colors

Update the CSS custom properties in `src/app/globals.css` to change the base color palette:

```css
/* CanTrip Custom Colors */
--color-urban-blue: #3A5A8C;
--color-maple-red: #D32F2F;
--color-forest-green: #388E3C;
/* ... */
```

## Development

- **Linting**: ESLint is configured for code quality
- **Type Checking**: TypeScript ensures type safety
- **Hot Reload**: Changes are reflected immediately during development

## Production Build

```bash
npm run build
npm start
```

## License

This project is part of the CanTrip travel assistant platform.