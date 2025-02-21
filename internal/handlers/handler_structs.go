package handlers

type SteamInventoryResponse struct {
	Assets              []Asset       `json:"assets"`
	Descriptions        []Description `json:"descriptions"`
	MoreItems           int           `json:"more_items"`
	LastAssetid         string        `json:"last_assetid"`
	TotalInventoryCount int           `json:"total_inventory_count"`
	Success             int           `json:"success"`
	Rwgrsn              int           `json:"rwgrsn"`
}

type Asset struct {
	AppID      int    `json:"appid"`
	ContextID  string `json:"contextid"`
	AssetID    string `json:"assetid"`
	ClassID    string `json:"classid"`
	InstanceID string `json:"instanceid"`
	Amount     string `json:"amount"`
}

type Description struct {
	Appid                       int                `json:"appid"`
	Classid                     string             `json:"classid"`
	Instanceid                  string             `json:"instanceid"`
	Currency                    int                `json:"currency"`
	BackgroundColor             string             `json:"background_color"`
	IconURL                     string             `json:"icon_url"`
	IconURLLarge                string             `json:"icon_url_large"`
	Descriptions                []DescriptionValue `json:"descriptions"`
	Tradable                    int                `json:"tradable"`
	Name                        string             `json:"name"`
	Type                        string             `json:"type"`
	MarketName                  string             `json:"market_name"`
	MarketHashName              string             `json:"market_hash_name"`
	MarketFeeApp                int                `json:"market_fee_app"`
	Commodity                   int                `json:"commodity"`
	MarketTradableRestriction   int                `json:"market_tradable_restriction"`
	MarketMarketableRestriction int                `json:"market_marketable_restriction"`
	Marketable                  int                `json:"marketable"`
	Tags                        []Tag              `json:"tags"`
}

type DescriptionValue struct {
	Value string `json:"value"`
}

type Tag struct {
	Category              string `json:"category"`
	InternalName          string `json:"internal_name"`
	LocalizedCategoryName string `json:"localized_category_name"`
	LocalizedTagName      string `json:"localized_tag_name"`
}

type SteamTradeResponse struct {
	Response struct {
		Items []struct {
			AppID      int    `json:"appid"`
			ContextID  string `json:"contextid"`
			AssetID    string `json:"assetid"`
			ClassID    string `json:"classid"`
			InstanceID string `json:"instanceid"`
			Amount     string `json:"amount"`
			Desc       struct {
				Type       string `json:"type"`
				Name       string `json:"name"`
				MarketName string `json:"market_name"`
				MarketHash string `json:"market_hash_name"`
				IconURL    string `json:"icon_url"`
				Tradable   int    `json:"tradable"`
				Marketable int    `json:"marketable"`
			} `json:"descriptions"`
		} `json:"items"`
		TotalInventoryCount int `json:"total_inventory_count"`
	} `json:"response"`
}

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

type GameFromList struct {
	AppID                    int    `json:"appid"`
	Name                     string `json:"name"`
	PlaytimeForever          int    `json:"playtime_forever"`
	ImgIconURL               string `json:"img_icon_url"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats,omitempty"`
	PlaytimeWindowsForever   int    `json:"playtime_windows_forever"`
	PlaytimeMacForever       int    `json:"playtime_mac_forever"`
	PlaytimeLinuxForever     int    `json:"playtime_linux_forever"`
	PlaytimeDeckForever      int    `json:"playtime_deck_forever"`
	RtimeLastPlayed          int    `json:"rtime_last_played"`
	CapsuleFilename          string `json:"capsule_filename"`
	HasWorkshop              bool   `json:"has_workshop"`
	HasMarket                bool   `json:"has_market"`
	HasDlc                   bool   `json:"has_dlc"`
	ContentDescriptorids     []int  `json:"content_descriptorids,omitempty"`
	PlaytimeDisconnected     int    `json:"playtime_disconnected"`
	HasLeaderboards          bool   `json:"has_leaderboards,omitempty"`
	SortAs                   string `json:"sort_as,omitempty"`
	Playtime2Weeks           int    `json:"playtime_2weeks,omitempty"`
}

type GameDataFull struct {
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

type PlayerResponse struct {
	Response struct {
		Players []struct {
			Steamid                  string `json:"steamid"`
			Communityvisibilitystate int    `json:"communityvisibilitystate"`
			Profilestate             int    `json:"profilestate"`
			Personaname              string `json:"personaname"`
			Commentpermission        int    `json:"commentpermission"`
			Profileurl               string `json:"profileurl"`
			Avatar                   string `json:"avatar"`
			Avatarmedium             string `json:"avatarmedium"`
			Avatarfull               string `json:"avatarfull"`
			Avatarhash               string `json:"avatarhash"`
			Lastlogoff               int    `json:"lastlogoff"`
			Personastate             int    `json:"personastate"`
			Realname                 string `json:"realname"`
			Primaryclanid            string `json:"primaryclanid"`
			Timecreated              int    `json:"timecreated"`
			Personastateflags        int    `json:"personastateflags"`
		} `json:"players"`
	} `json:"response"`
}
