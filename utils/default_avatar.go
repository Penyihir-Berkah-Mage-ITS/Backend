package utils

func DefaultAvatar(input string) string {
	avatarLinks := map[string]string{
		"1": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/1.png",
		"2": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/2.png",
		"3": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/3.png",
		"4": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/4.png",
		"5": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/5.png",
		"6": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/6.png",
		"7": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/7.png",
		"8": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/8.png",
		"9": "https://pbhlfzbgpnywavonhdld.supabase.co/storage/v1/object/public/images/9.png",
	}

	if link, exists := avatarLinks[input]; exists {
		return link
	}

	return "https://defaultlink.com"
}
