package filter

// DefaultFilter initializes a new Filter with predefined patterns
func DefaultFilter() *Filter {
	return &Filter{
		ExcludePatterns: []string{
			// "Data/enUS/Documentation/*",   // Don't ignore documentation
		},
		ExactMatches: []string{
			"README.md",
			// Add more exact filenames as needed
		},
		ExtensionMatches: []string{
			".gitignore",
			".env",
			".log",
			// Add more file extensions to ignore
		},
		BaseMatches: []string{
			"manifest.json",
			"Battle.net.dll",
			// base ChromieCraft 3.3.5a files
			"dbghelp.dll",
			"DivxDecoder.dll",
			"ijl15.dll",
			"msvcr80.dll",
			"Repair.exe",
			"Scan.dll",
			"unicows.dll",
			"WowError.exe",
			"Wow.exe",
			"Data/common-2.MPQ",
			"Data/common.MPQ",
			"Data/expansion.MPQ",
			"Data/lichking.MPQ",
			"Data/patch-2.MPQ",
			"Data/patch-3.MPQ",
			"Data/patch.MPQ",
			"Data/enUS/AccountBilling.url",
			"Data/enUS/backup-enUS.MPQ",
			"Data/enUS/base-enUS.MPQ",
			"Data/enUS/connection-help.html",
			"Data/enUS/Credits_BC.html",
			"Data/enUS/Credits.html",
			"Data/enUS/Credits_LK.html",
			"Data/enUS/eula.html",
			"Data/enUS/expansion-locale-enUS.MPQ",
			"Data/enUS/expansion-speech-enUS.MPQ",
			"Data/enUS/lichking-locale-enUS.MPQ",
			"Data/enUS/lichking-speech-enUS.MPQ",
			"Data/enUS/locale-enUS.MPQ",
			"Data/enUS/patch-enUS-2.MPQ",
			"Data/enUS/patch-enUS-3.MPQ",
			"Data/enUS/patch-enUS.MPQ",
			"Data/enUS/realmlist.wtf",
			"Data/enUS/speech-enUS.MPQ",
			"Data/enUS/TechSupport.url",
			"Data/enUS/tos.html",
			// non-base files
			"WoW.exe",
			// Add more base files as needed
		},
		GlobPatterns: []string{
			"patcher*",
			".wine/*",
			// base ChromieCraft 3.3.5a files
			"Data/enUS/Interface/Cinematics/*.avi", // Ignore all cinematic files
			"Data/enUS/Documentation/*",            // Ignore all documentation
			"Interface/AddOns/Blizzard_*",          // Ignore all Blizzard addons
			// non-base files
			"Interface/AddOns/*", // Ignore all user addons
			"WTF/*",              // Ignore all Warcraft Text Files (WTF)
			"Cache/*",            // Ignore cache
			"Logs/*",             // Ignore logs
			"Errors/*",           // Ignore errors
			"Screenshots/*",      // Ignore screenshots
			// Add more glob patterns as needed
		},
	}
}
