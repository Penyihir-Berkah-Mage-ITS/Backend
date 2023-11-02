package utils

func DefaultAvatar(input string) string {
	avatarLinks := map[string]string{
		"1": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_1.png",
		"2": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_2.png",
		"3": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_3.png",
		"4": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_4.png",
		"5": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_5.png",
		"6": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_6.png",
		"7": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_7.png",
		"8": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_8.png",
		"9": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/avatar_9.png",
	}

	if link, exists := avatarLinks[input]; exists {
		return link
	}

	return "https://defaultlink.com"
}
