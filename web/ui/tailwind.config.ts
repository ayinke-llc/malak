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
  			foreground: "hsl(222 47% 11%)",
  			card: {
  				DEFAULT: "hsl(0 0% 100%)",
  				foreground: "hsl(222 47% 11%)",
  			},
  			popover: {
  				DEFAULT: "hsl(0 0% 100%)",
  				foreground: "hsl(222 47% 11%)",
  			},
  			primary: {
  				DEFAULT: "hsl(160 84% 39%)",
  				foreground: "hsl(0 0% 100%)",
  			},
  			secondary: {
  				DEFAULT: "hsl(210 40% 96.1%)",
  				foreground: "hsl(222 47% 11%)",
  			},
  			muted: {
  				DEFAULT: "hsl(210 40% 96.1%)",
  				foreground: "hsl(215.4 16.3% 46.9%)",
  			},
  			accent: {
  				DEFAULT: "hsl(210 40% 96.1%)",
  				foreground: "hsl(222 47% 11%)",
  			},
  			destructive: {
  				DEFAULT: "hsl(0 84.2% 60.2%)",
  				foreground: "hsl(210 40% 98%)",
  			},
  			border: "hsl(214.3 31.8% 91.4%)",
  			input: "hsl(214.3 31.8% 91.4%)",
  			ring: "hsl(222 47% 11%)",
  			chart: {
  				"1": "hsl(222 47% 11%)",
  				"2": "hsl(215.4 16.3% 46.9%)",
  				"3": "hsl(214.3 31.8% 91.4%)",
  				"4": "hsl(210 40% 96.1%)",
  				"5": "hsl(0 0% 100%)",
  			},
  			sidebar: {
  				DEFAULT: "hsl(0 0% 100%)",
  				foreground: "hsl(222 47% 11%)",
  				primary: "hsl(222 47% 11%)",
  				"primary-foreground": "hsl(0 0% 100%)",
  				accent: "hsl(210 40% 96.1%)",
  				"accent-foreground": "hsl(222 47% 11%)",
  				border: "hsl(214.3 31.8% 91.4%)",
  				ring: "hsl(222 47% 11%)",
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
