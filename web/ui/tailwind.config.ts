import type { Config } from "tailwindcss";

const config: Config = {
  darkMode: ["selector", "class"],
  content: ["./src/**/*.{js,ts,jsx,tsx,mdx}"],
  theme: {
  	extend: {
  		keyframes: {
  			hide: {
  				from: {
  					opacity: '1'
  				},
  				to: {
  					opacity: '0'
  				}
  			},
  			slideDownAndFade: {
  				from: {
  					opacity: '0',
  					transform: 'translateY(-6px)'
  				},
  				to: {
  					opacity: '1',
  					transform: 'translateY(0)'
  				}
  			},
  			slideLeftAndFade: {
  				from: {
  					opacity: '0',
  					transform: 'translateX(6px)'
  				},
  				to: {
  					opacity: '1',
  					transform: 'translateX(0)'
  				}
  			},
  			slideUpAndFade: {
  				from: {
  					opacity: '0',
  					transform: 'translateY(6px)'
  				},
  				to: {
  					opacity: '1',
  					transform: 'translateY(0)'
  				}
  			},
  			slideRightAndFade: {
  				from: {
  					opacity: '0',
  					transform: 'translateX(-6px)'
  				},
  				to: {
  					opacity: '1',
  					transform: 'translateX(0)'
  				}
  			},
  			dialogOverlayShow: {
  				from: {
  					opacity: '0'
  				},
  				to: {
  					opacity: '1'
  				}
  			},
  			dialogContentShow: {
  				from: {
  					opacity: '0',
  					transform: 'translate(-50%, -45%) scale(0.95)'
  				},
  				to: {
  					opacity: '1',
  					transform: 'translate(-50%, -50%) scale(1)'
  				}
  			},
  			drawerSlideLeftAndFade: {
  				from: {
  					opacity: '0',
  					transform: 'translateX(50%)'
  				},
  				to: {
  					opacity: '1',
  					transform: 'translateX(0)'
  				}
  			},
  			'accordion-down': {
  				from: {
  					height: '0'
  				},
  				to: {
  					height: 'var(--radix-accordion-content-height)'
  				}
  			},
  			'accordion-up': {
  				from: {
  					height: 'var(--radix-accordion-content-height)'
  				},
  				to: {
  					height: '0'
  				}
  			}
  		},
  		animation: {
  			hide: 'hide 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			slideDownAndFade: 'slideDownAndFade 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			slideLeftAndFade: 'slideLeftAndFade 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			slideUpAndFade: 'slideUpAndFade 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			slideRightAndFade: 'slideRightAndFade 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			drawerSlideLeftAndFade: 'drawerSlideLeftAndFade 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			dialogOverlayShow: 'dialogOverlayShow 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			dialogContentShow: 'dialogContentShow 150ms cubic-bezier(0.16, 1, 0.3, 1)',
  			'accordion-down': 'accordion-down 0.2s ease-out',
  			'accordion-up': 'accordion-up 0.2s ease-out'
  		},
  		borderRadius: {
  			lg: 'var(--radius)',
  			md: 'calc(var(--radius) - 2px)',
  			sm: 'calc(var(--radius) - 4px)'
  		},
  		colors: {
  			background: "hsl(0 0% 100%)",
  			foreground: "#2c4975",
  			card: {
  				DEFAULT: "hsl(0 0% 100%)",
  				foreground: "#2c4975",
  			},
  			popover: {
  				DEFAULT: "hsl(0 0% 100%)",
  				foreground: "#2c4975",
  			},
  			primary: {
  				DEFAULT: "#0a3171",
  				foreground: "#ffffff",
  			},
  			secondary: {
  				DEFAULT: "#1150db",
  				foreground: "#ffffff",
  			},
  			tertiary: {
  				DEFAULT: "#4178cd",
  				foreground: "#ffffff",
  			},
  			muted: {
  				DEFAULT: "#e6edf7",
  				foreground: "#2c4975",
  			},
  			accent: {
  				DEFAULT: "#1150db",
  				foreground: "#ffffff",
  			},
  			destructive: {
  				DEFAULT: "hsl(0 84.2% 60.2%)",
  				foreground: "#ffffff",
  			},
  			border: "#e6edf7",
  			input: "#e6edf7",
  			ring: "#0a3171",
  			theme: {
  				light: "#f0f5fc",
  				dark: "#092657"
  			},
  			body: "#ffffff",
  			text: {
  				DEFAULT: "#2c4975",
  				dark: "#0a3171",
  				light: "#f0f5fc"
  			}
  		}
  	}
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
    require("tailwindcss-animate"),
  ],
};

export default config;
