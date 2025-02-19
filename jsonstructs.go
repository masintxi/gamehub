package main

type GameData struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string `json:"type"`
		Name                string `json:"name"`
		SteamAppid          int    `json:"steam_appid"`
		IsFree              bool   `json:"is_free"`
		ControllerSupport   string `json:"controller_support"`
		Dlc                 []int  `json:"dlc"`
		DetailedDescription string `json:"detailed_description"`
		AboutTheGame        string `json:"about_the_game"`
		ShortDescription    string `json:"short_description"`
		SupportedLanguages  string `json:"supported_languages"`
		PriceOverview       struct {
			Currency         string `json:"currency"`
			Initial          int    `json:"initial"`
			Final            int    `json:"final"`
			DiscountPercent  int    `json:"discount_percent"`
			InitialFormatted string `json:"initial_formatted"`
			FinalFormatted   string `json:"final_formatted"`
		} `json:"price_overview"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
	} `json:"data"`
}

type GameDataold struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string `json:"type"`
		Name                string `json:"name"`
		SteamAppid          int    `json:"steam_appid"`
		RequiredAge         int    `json:"required_age"`
		IsFree              bool   `json:"is_free"`
		ControllerSupport   string `json:"controller_support"`
		Dlc                 []int  `json:"dlc"`
		DetailedDescription string `json:"detailed_description"`
		AboutTheGame        string `json:"about_the_game"`
		ShortDescription    string `json:"short_description"`
		SupportedLanguages  string `json:"supported_languages"`
		HeaderImage         string `json:"header_image"`
		CapsuleImage        string `json:"capsule_image"`
		CapsuleImagev5      string `json:"capsule_imagev5"`
		Website             string `json:"website"`
		PcRequirements      struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"pc_requirements"`
		MacRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"mac_requirements"`
		LinuxRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"linux_requirements"`
		LegalNotice   string   `json:"legal_notice"`
		Developers    []string `json:"developers"`
		Publishers    []string `json:"publishers"`
		PriceOverview struct {
			Currency         string `json:"currency"`
			Initial          int    `json:"initial"`
			Final            int    `json:"final"`
			DiscountPercent  int    `json:"discount_percent"`
			InitialFormatted string `json:"initial_formatted"`
			FinalFormatted   string `json:"final_formatted"`
		} `json:"price_overview"`
		Packages      []int `json:"packages"`
		PackageGroups []struct {
			Name                    string `json:"name"`
			Title                   string `json:"title"`
			Description             string `json:"description"`
			SelectionText           string `json:"selection_text"`
			SaveText                string `json:"save_text"`
			DisplayType             int    `json:"display_type"`
			IsRecurringSubscription string `json:"is_recurring_subscription"`
			Subs                    []struct {
				Packageid                int    `json:"packageid"`
				PercentSavingsText       string `json:"percent_savings_text"`
				PercentSavings           int    `json:"percent_savings"`
				OptionText               string `json:"option_text"`
				OptionDescription        string `json:"option_description"`
				CanGetFreeLicense        string `json:"can_get_free_license"`
				IsFreeLicense            bool   `json:"is_free_license"`
				PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
			} `json:"subs"`
		} `json:"package_groups"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"genres"`
		Screenshots []struct {
			ID            int    `json:"id"`
			PathThumbnail string `json:"path_thumbnail"`
			PathFull      string `json:"path_full"`
		} `json:"screenshots"`
		Movies []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm      struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Mp4 struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"mp4"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements struct {
			Total       int `json:"total"`
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted"`
		} `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background         string `json:"background"`
		BackgroundRaw      string `json:"background_raw"`
		ContentDescriptors struct {
			Ids   []int  `json:"ids"`
			Notes string `json:"notes"`
		} `json:"content_descriptors"`
		Ratings struct {
			Esrb struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"esrb"`
			Oflc struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"oflc"`
			Nzoflc struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"nzoflc"`
			Cero struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"cero"`
			Pegi struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"pegi"`
			Usk struct {
				Rating      string `json:"rating"`
				RatingID    string `json:"rating_id"`
				Descriptors string `json:"descriptors"`
			} `json:"usk"`
			Kgrb struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"kgrb"`
			Dejus struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"dejus"`
			Fpb struct {
				Rating string `json:"rating"`
			} `json:"fpb"`
			Csrr struct {
				Rating      string `json:"rating"`
				Descriptors string `json:"descriptors"`
			} `json:"csrr"`
			SteamGermany struct {
				RatingGenerated string `json:"rating_generated"`
				Rating          string `json:"rating"`
				RequiredAge     string `json:"required_age"`
				Banned          string `json:"banned"`
				UseAgeGate      string `json:"use_age_gate"`
				Descriptors     string `json:"descriptors"`
			} `json:"steam_germany"`
		} `json:"ratings"`
	} `json:"data"`
}
