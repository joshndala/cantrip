export interface CityTheme {
  primary: string;
  secondary: string;
  accent: string;
  background: string;
  surface: string;
  text: string;
  gradient: string;
}

export const cityThemes: Record<string, CityTheme> = {
  // Major Cities
  toronto: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#1E88E5', // Lake Blue
    accent: '#E0E0E0', // Concrete Grey
    background: 'linear-gradient(135deg, #3A5A8C 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-urban-blue via-lake-blue to-concrete-grey'
  },
  
  vancouver: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#388E3C', // Forest Green
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #388E3C 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-lake-blue to-sunrise-gold'
  },
  
  montreal: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#D32F2F', // Maple Red
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #D32F2F 0%, #FFC107 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-maple-red via-sunrise-gold to-forest-green'
  },
  
  calgary: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#FFC107', // Sunrise Gold
    accent: '#D32F2F', // Maple Red
    background: 'linear-gradient(135deg, #FFC107 0%, #D32F2F 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-sunrise-gold via-maple-red to-concrete-grey'
  },
  
  ottawa: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#D32F2F', // Maple Red
    accent: '#E0E0E0', // Concrete Grey
    background: 'linear-gradient(135deg, #3A5A8C 0%, #D32F2F 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-urban-blue via-maple-red to-concrete-grey'
  },
  
  edmonton: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#388E3C', // Forest Green
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #388E3C 0%, #FFC107 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-sunrise-gold to-concrete-grey'
  },
  
  halifax: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#1E88E5', // Lake Blue
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #1E88E5 0%, #3A5A8C 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-lake-blue via-urban-blue to-sunrise-gold'
  },
  
  'kitchener-waterloo': {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#E0E0E0', // Concrete Grey
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #3A5A8C 0%, #E0E0E0 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-urban-blue via-concrete-grey to-sunrise-gold'
  },
  
  gatineau: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#388E3C', // Forest Green
    accent: '#D32F2F', // Maple Red
    background: 'linear-gradient(135deg, #388E3C 0%, #D32F2F 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-maple-red to-concrete-grey'
  },
  
  // Provincial Capitals
  'quebec-city': {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#D32F2F', // Maple Red
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #D32F2F 0%, #FFC107 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-maple-red via-sunrise-gold to-forest-green'
  },
  
  victoria: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#388E3C', // Forest Green
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #388E3C 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-lake-blue to-sunrise-gold'
  },
  
  // Resort & Tourist Destinations
  banff: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#388E3C', // Forest Green
    accent: '#FAFAFA', // Snow White
    background: 'linear-gradient(135deg, #388E3C 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-lake-blue to-snow-white'
  },
  
  whistler: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#1E88E5', // Lake Blue
    accent: '#E0E0E0', // Concrete Grey
    background: 'linear-gradient(135deg, #3A5A8C 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-urban-blue via-lake-blue to-concrete-grey'
  },
  
  jasper: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#388E3C', // Forest Green
    accent: '#E0E0E0', // Concrete Grey
    background: 'linear-gradient(135deg, #388E3C 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-lake-blue to-concrete-grey'
  },
  
  'niagara-region': {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#1E88E5', // Lake Blue
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #1E88E5 0%, #FFC107 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-lake-blue via-sunrise-gold to-maple-red'
  },
  
  // Specialized Destinations
  yukon: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#1B263B', // Midnight Blue
    accent: '#FAFAFA', // Snow White
    background: 'linear-gradient(135deg, #1B263B 0%, #388E3C 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-midnight-blue via-forest-green to-snow-white'
  },
  
  'gros-morne': {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#1565C0', // Ocean Blue
    accent: '#E0E0E0', // Concrete Grey
    background: 'linear-gradient(135deg, #388E3C 0%, #1565C0 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-ocean-blue to-concrete-grey'
  },
  
  churchill: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#F5F5F5', // Ice White
    accent: '#1B263B', // Midnight Blue
    background: 'linear-gradient(135deg, #F5F5F5 0%, #1E88E5 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-ice-white via-lake-blue to-midnight-blue'
  },
  
  'cape-breton': {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#1565C0', // Ocean Blue
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #1565C0 0%, #388E3C 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-ocean-blue via-forest-green to-sunrise-gold'
  },
  
  saguenay: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#388E3C', // Forest Green
    accent: '#D32F2F', // Maple Red
    background: 'linear-gradient(135deg, #388E3C 0%, #D32F2F 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-forest-green via-maple-red to-concrete-grey'
  },
  
  kingston: {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#E0E0E0', // Concrete Grey
    accent: '#FFC107', // Sunrise Gold
    background: 'linear-gradient(135deg, #3A5A8C 0%, #E0E0E0 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-urban-blue via-concrete-grey to-sunrise-gold'
  },
  
  'trois-rivieres': {
    primary: '#3A5A8C', // Urban Blue
    secondary: '#D32F2F', // Maple Red
    accent: '#E0E0E0', // Concrete Grey
    background: 'linear-gradient(135deg, #D32F2F 0%, #388E3C 100%)',
    surface: '#F5F5F5',
    text: '#1A1A1A',
    gradient: 'from-maple-red via-concrete-grey to-forest-green'
  }
};

// Default neutral theme
export const defaultTheme: CityTheme = {
  primary: '#3A5A8C', // Urban Blue
  secondary: '#E0E0E0', // Light Grey
  accent: '#FFFFFF', // White
  background: 'linear-gradient(135deg, #E0E0E0 0%, #FFFFFF 100%)',
  surface: '#F5F5F5',
  text: '#1A1A1A',
  gradient: 'from-light-grey to-white'
};

// Helper function to get theme by city name
export function getCityTheme(cityName: string): CityTheme {
  const normalizedCity = cityName.toLowerCase().replace(/\s+/g, '-');
  return cityThemes[normalizedCity] || defaultTheme;
}

// Helper function to detect city from message
export function detectCityFromMessage(message: string): string | null {
  const cityKeywords = Object.keys(cityThemes).map(city => 
    city.replace('-', ' ').replace('-', ' ')
  );
  
  const lowerMessage = message.toLowerCase();
  
  for (const city of cityKeywords) {
    if (lowerMessage.includes(city)) {
      return city.replace(/\s+/g, '-');
    }
  }
  
  // Check for common city variations
  const cityVariations: Record<string, string> = {
    'toronto': 'toronto',
    'vancouver': 'vancouver',
    'montreal': 'montreal',
    'calgary': 'calgary',
    'ottawa': 'ottawa',
    'edmonton': 'edmonton',
    'halifax': 'halifax',
    'quebec city': 'quebec-city',
    'quebec': 'quebec-city',
    'victoria': 'victoria',
    'banff': 'banff',
    'whistler': 'whistler',
    'jasper': 'jasper',
    'niagara': 'niagara-region',
    'yukon': 'yukon',
    'gros morne': 'gros-morne',
    'churchill': 'churchill',
    'cape breton': 'cape-breton',
    'saguenay': 'saguenay',
    'kingston': 'kingston',
    'trois-rivières': 'trois-rivieres',
    'trois rivieres': 'trois-rivieres'
  };
  
  for (const [keyword, cityKey] of Object.entries(cityVariations)) {
    if (lowerMessage.includes(keyword)) {
      return cityKey;
    }
  }
  
  return null;
}
