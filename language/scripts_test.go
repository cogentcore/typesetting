package language

import (
	"testing"
	"unicode"

	tu "github.com/cogentcore/typesetting/testutils"
)

func TestParseScript(t *testing.T) {
	tests := []struct {
		args    string
		want    Script
		wantErr bool
	}{
		{"xxx", 0, true},
		{"bamu", Bamum, false},
		{"bamu_to_long", Bamum, false},
		{"cyrl", Cyrillic, false},
		{"samr", Samaritan, false},
		{"ARAB", Arabic, false},
		{"arab", Arabic, false},
		{"Arab", Arabic, false},
		{"Samr", Samaritan, false},
	}
	for _, tt := range tests {
		got, err := ParseScript(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseScript() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("ParseScript() = %v, want %v", got, tt.want)
		}
	}
}

func TestLookupScript(t *testing.T) {
	for _, test := range []struct {
		r      rune
		script Script
	}{
		{'a', Latin},
		{'ل', Arabic},
		{0, Common},
		{0xFFFFFFF, Unknown},
		// copied from harfbuzz
		{0x002A, Common},
		{0x0670, Inherited},
		{0x060D, Arabic},
		{0x0559, Armenian},
		{0x09CD, Bengali},
		{0x31B6, Bopomofo},
		{0x13A2, Cherokee},
		{0x2CFD, Coptic},
		{0x0482, Cyrillic},
		{0x10401, Deseret},
		{0x094D, Devanagari},
		{0x1258, Ethiopic},
		{0x10FC, Georgian},
		{0x10341, Gothic},
		{0x0375, Greek},
		{0x0A83, Gujarati},
		{0x0A3C, Gurmukhi},
		{0x3005, Han},
		{0x1100, Hangul},
		{0x05BF, Hebrew},
		{0x309F, Hiragana},
		{0x0CBC, Kannada},
		{0x30FF, Katakana},
		{0x17DD, Khmer},
		{0x0EDD, Lao},
		{0x0061, Latin},
		{0x0D3D, Malayalam},
		{0x1843, Mongolian},
		{0x1031, Myanmar},
		{0x169C, Ogham},
		{0x10322, Old_Italic},
		{0x0B3C, Oriya},
		{0x16EF, Runic},
		{0x0DBD, Sinhala},
		{0x0711, Syriac},
		{0x0B82, Tamil},
		{0x0C03, Telugu},
		{0x07B1, Thaana},
		{0x0E31, Thai},
		{0x0FD4, Tibetan},
		// {0x1401, Canadian_Syllabics},
		{0xA015, Yi},
		{0x1700, Tagalog},
		{0x1720, Hanunoo},
		{0x1740, Buhid},
		{0x1760, Tagbanwa},

		/* Unicode-4.0 additions */
		{0x2800, Braille},
		{0x10808, Cypriot},
		{0x1932, Limbu},
		{0x10480, Osmanya},
		{0x10450, Shavian},
		{0x10000, Linear_B},
		{0x1950, Tai_Le},
		{0x1039F, Ugaritic},

		/* Unicode-4.1 additions */
		{0x1980, New_Tai_Lue},
		{0x1A1F, Buginese},
		{0x2C00, Glagolitic},
		{0x2D6F, Tifinagh},
		{0xA800, Syloti_Nagri},
		{0x103D0, Old_Persian},
		{0x10A3F, Kharoshthi},

		/* Unicode-5.0 additions */
		{0x0378, Unknown},
		{0x1B04, Balinese},
		{0x12000, Cuneiform},
		{0x10900, Phoenician},
		{0xA840, Phags_Pa},
		{0x07C0, Nko},

		/* Unicode-5.1 additions */
		{0xA900, Kayah_Li},
		{0x1C00, Lepcha},
		{0xA930, Rejang},
		{0x1B80, Sundanese},
		{0xA880, Saurashtra},
		{0xAA00, Cham},
		{0x1C50, Ol_Chiki},
		{0xA500, Vai},
		{0x102A0, Carian},
		{0x10280, Lycian},
		{0x1093F, Lydian},

		{0x111111, Unknown},

		/* Unicode-5.2 additions */
		{0x10B00, Avestan},
		{0xA6A0, Bamum},
		{0x1400, Canadian_Aboriginal},
		{0x13000, Egyptian_Hieroglyphs},
		{0x10840, Imperial_Aramaic},
		{0x1CED, Inherited},
		{0x10B60, Inscriptional_Pahlavi},
		{0x10B40, Inscriptional_Parthian},
		{0xA980, Javanese},
		{0x11082, Kaithi},
		{0xA4D0, Lisu},
		{0xABE5, Meetei_Mayek},
		{0x10A60, Old_South_Arabian},
		{0x10C00, Old_Turkic},
		{0x0800, Samaritan},
		{0x1A20, Tai_Tham},
		{0xAA80, Tai_Viet},

		/* Unicode-6.0 additions */
		{0x1BC0, Batak},
		{0x11000, Brahmi},
		{0x0840, Mandaic},

		/* Unicode-6.1 additions */
		{0x10980, Meroitic_Hieroglyphs},
		{0x109A0, Meroitic_Cursive},
		{0x110D0, Sora_Sompeng},
		{0x11100, Chakma},
		{0x11180, Sharada},
		{0x11680, Takri},
		{0x16F00, Miao},

		/* Unicode-6.2 additions */
		{0x20BA, Common},

		/* Unicode-6.3 additions */
		{0x2066, Common},

		/* Unicode-7.0 additions */
		{0x10350, Old_Permic},
		{0x10500, Elbasan},
		{0x10530, Caucasian_Albanian},
		{0x10600, Linear_A},
		{0x10860, Palmyrene},
		{0x10880, Nabataean},
		{0x10A80, Old_North_Arabian},
		{0x10AC0, Manichaean},
		{0x10B80, Psalter_Pahlavi},
		{0x11150, Mahajani},
		{0x11200, Khojki},
		{0x112B0, Khudawadi},
		{0x11300, Grantha},
		{0x11480, Tirhuta},
		{0x11580, Siddham},
		{0x11600, Modi},
		{0x118A0, Warang_Citi},
		{0x11AC0, Pau_Cin_Hau},
		{0x16A40, Mro},
		{0x16AD0, Bassa_Vah},
		{0x16B00, Pahawh_Hmong},
		{0x1BC00, Duployan},
		{0x1E800, Mende_Kikakui},

		/* Unicode-8.0 additions */
		{0x108E0, Hatran},
		{0x10C80, Old_Hungarian},
		{0x11280, Multani},
		{0x11700, Ahom},
		{0x14400, Anatolian_Hieroglyphs},
		{0x1D800, SignWriting},

		/* Unicode-9.0 additions */
		{0x104B0, Osage},
		{0x11400, Newa},
		{0x11C00, Bhaiksuki},
		{0x11C70, Marchen},
		{0x17000, Tangut},
		{0x1E900, Adlam},

		/* Unicode-10.0 additions */
		{0x11A00, Zanabazar_Square},
		{0x11A50, Soyombo},
		{0x11D00, Masaram_Gondi},
		{0x1B170, Nushu},

		/* Unicode-11.0 additions */
		{0x10D00, Hanifi_Rohingya},
		{0x10F00, Old_Sogdian},
		{0x10F30, Sogdian},
		{0x11800, Dogra},
		{0x11D60, Gunjala_Gondi},
		{0x11EE0, Makasar},
		{0x16E40, Medefaidrin},

		/* Unicode-12.0 additions */
		{0x10FE0, Elymaic},
		{0x119A0, Nandinagari},
		{0x1E100, Nyiakeng_Puachue_Hmong},
		{0x1E2C0, Wancho},

		/* Unicode-12.1 additions */
		{0x32FF, Common},

		/* Unicode-13.0 additions */
		{0x10E80, Yezidi},
		{0x10FB0, Chorasmian},
		{0x11900, Dives_Akuru},
		{0x18B00, Khitan_Small_Script},

		/* Unicode-14.0 additions */
		{0x10570, Vithkuqi},
		{0x10F70, Old_Uyghur},
		{0x12F90, Cypro_Minoan},
		{0x16A70, Tangsa},
		{0x1E290, Toto},

		/* Unicode-15.0 additions */
		{0x11F00, Kawi},
		{0x1E4D0, Nag_Mundari},

		/* Unicode-16.0 additions */
		{0x105C0, Todhri},
		{0x10D40, Garay},
		{0x11380, Tulu_Tigalari},
		{0x11BC0, Sunuwar},
		{0x16100, Gurung_Khema},
		{0x16D40, Kirat_Rai},
		{0x1E5D0, Ol_Onal},

		/* Unicode-16.0 additions */
		{0x10940, Sidetic},
		{0x11DB0, Tolong_Siki},
		{0x16EA0, Beria_Erfe},
		{0x1E6C0, Tai_Yo},

		{0x111111, Unknown},
	} {
		tu.Assert(t, LookupScript(test.r) == test.script)
	}
}

func TestScript_String(t *testing.T) {
	tu.Assert(t, Bamum.String() == "Bamu")
}

func TestScript_Strong(t *testing.T) {
	tu.Assert(t, Latin.Strong())
	tu.Assert(t, Arabic.Strong())
	tu.Assert(t, Unknown.Strong())
	tu.Assert(t, !Common.Strong())
	tu.Assert(t, !Inherited.Strong())
}

// used as benchmark reference
func lookupScriptNaive(r rune) Script {
	for name, table := range unicode.Scripts {
		if unicode.Is(table, r) {
			return scriptToTag[name]
		}
	}
	return Unknown
}

func TestFastLookup(t *testing.T) {
	for _, r := range scriptsSample {
		g1, g2 := lookupScriptNaive(r), LookupScript(r)
		if g1 != g2 {
			t.Fatalf("for rune 0x%x, expected %s, got %s", r, g1, g2)
		}
	}
}

func BenchmarkLookupScript(b *testing.B) {
	b.Run("naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, r := range scriptsSample {
				_ = lookupScriptNaive(r)
			}
		}
	})
	b.Run("optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, r := range scriptsSample {
				_ = LookupScript(r)
			}
		}
	})
}

//lint:ignore ST1018 for simplicity
const scriptsSample = `
	Ek kan glas eet, maar dit doen my nie skade nie. 
	نص حكيم له سر قاطع وذو شأن عظيم مكتوب على ثوب أخضر ومغلف بجلد أزرق. 
	Gvxam mincetu apocikvyeh: ñizol ce mamvj ka raq kuse bafkeh mew. 
	I koh Glos esa, und es duard ma ned wei. 
	Под южно дърво, цъфтящо в синьо, бягаше малко пухкаво зайче. 
	Mi save kakae glas, hemi no save katem mi. 
	আমি কাঁচ খেতে পারি, তাতে আমার কোনো ক্ষতি হয় না। 
	ཤེལ་སྒོ་ཟ་ནས་ང་ན་གི་མ་རེད། 
	Fin džip, gluh jež i čvrst konjić dođoše bez moljca. 
	Jove xef, porti whisky amb quinze glaçons d'hidrogen, coi! 
	Siña yo' chumocho krestat, ti ha na'lalamen yo'. 
	Příliš žluťoučký kůň úpěl ďábelské ódy. 
	Dw i'n gallu bwyta gwydr, 'dyw e ddim yn gwneud dolur i mi. 
	Quizdeltagerne spiste jordbær med fløde, mens cirkusklovnen Walther spillede på xylofon. 
	Zwölf Boxkämpfer jagen Viktor quer über den großen Sylter Deich. 
	މާއްދާ 1 – ހުރިހާ އިންސާނުން ވެސް އުފަންވަނީ، ދަރަޖަ އާއި ޙައްޤު ތަކުގައި މިނިވަންކަމާއި ހަމަހަމަކަން ލިބިގެންވާ ބައެއްގެ ގޮތުގައެވެ. 
	Θέλει αρετή και τόλμη η ελευθερία. (Ανδρέας Κάλβος) 
	The quick brown fox jumps over the lazy dog.
	Ich canne glas eten and hit hirtiþ me nouȝt. 
	Eĥoŝanĝo ĉiuĵaŭde. 
	Jovencillo emponzoñado de whisky: ¡qué figurota exhibe! 
	See väike mölder jõuab rongile hüpata. 
	Kristala jan dezaket, ez dit minik ematen. 
	«الا یا اَیُّها السّاقی! اَدِرْ کَأساً وَ ناوِلْها!» که عشق آسان نمود اوّل، ولی افتاد مشکل‌ها!
	Viekas kettu punaturkki laiskan koiran takaa kurkki. 
	Voix ambiguë d'un cœur qui, au zéphyr, préfère les jattes de kiwis. 
	Je puis mangier del voirre. Ne me nuit. 
	Chuaigh bé mhórshách le dlúthspád fíorfhinn trí hata mo dhea-phorcáin bhig. 
	S urrainn dhomh gloinne ithe; cha ghoirtich i mi. 
	Eu podo xantar cristais e non cortarme. 
	𐌼𐌰𐌲 𐌲𐌻𐌴𐍃 𐌹̈𐍄𐌰𐌽, 𐌽𐌹 𐌼𐌹𐍃 𐍅𐌿 𐌽𐌳𐌰𐌽 𐌱𐍂𐌹𐌲𐌲𐌹𐌸. 
	હું કાચ ખાઇ શકુ છુ અને તેનાથી મને દર્દ નથી થતુ. 
	Foddym gee glonney agh cha jean eh gortaghey mee. 
	Hiki iaʻu ke ʻai i ke aniani; ʻaʻole nō lā au e ʻeha. 
	דג סקרן שט לו בים זך אך לפתע פגש חבורה נחמדה שצצה כך. 
	नहीं नजर किसी की बुरी नहीं किसी का मुँह काला जो करे सो उपर वाला 
	Deblji krojač: zgužvah smeđ filc u tanjušni džepić. 
	Egy hűtlen vejét fülöncsípő, dühös mexikói úr Wesselényinél mázol Quitóban. 
	Կրնամ ապակի ուտել և ինծի անհանգիստ չըներ։ 
	Kæmi ný öxi hér ykist þjófum nú bæði víl og ádrepa 
	Ma la volpe, col suo balzo, ha raggiunto il quieto Fido. 
	いろはにほへと ちりぬるを 色は匂へど 散りぬるを
	Chruu, a kwik di kwik brong fox a jomp huova di liezi daag de, yu no siit?
	.o'i mu xagji sofybakni cu zvati le purdi 
	Aku isa mangan beling tanpa lara. 
	მინას ვჭამ და არა მტკივა. 
	ខ្ញុំអាចញុំកញ្ចក់បាន ដោយគ្មានបញ្ហារ 

	ನಾನು ಗಾಜನ್ನು ತಿನ್ನಬಲ್ಲೆ ಮತ್ತು ಅದರಿಂದ ನನಗೆ ನೋವಾಗುವುದಿಲ್ಲ. 
	다람쥐 헌 쳇바퀴에 타고파 
	Mý a yl dybry gwéder hag éf ny wra ow ankenya. 
	Sic surgens, dux, zelotypos quam karus haberis
	ຂອ້ຍກິນແກ້ວໄດ້ໂດຍທີ່ມັນບໍ່ໄດ້ເຮັດໃຫ້ຂອ້ຍເຈັບ 

	Įlinkdama fechtuotojo špaga sublykčiojusi pragręžė apvalų arbūzą. 
	Sarkanās jūrascūciņas peld pa jūru. 
	E koʻana e kai i te karahi, mea ʻā, ʻaʻe hauhau. 
	Можам да јадам стакло, а не ме штета. 
	വേദനയില്ലാതെ കുപ്പിചില്ലു് എനിയ്ക്കു് കഴിയ്ക്കാം. 
	ᠪᠢ ᠰᠢᠯᠢ ᠢᠳᠡᠶᠦ ᠴᠢᠳᠠᠨᠠ ᠂ ᠨᠠᠳᠤᠷ ᠬᠣᠤᠷᠠᠳᠠᠢ ᠪᠢᠰᠢ 
	मी काच खाऊ शकतो, मला ते दुखत नाही. 
	Saya boleh makan kaca dan ia tidak mencederakan saya. 
	ဘာသာပြန်နှင့် စာပေပြုစုရေး ကော်မရှင် 
	M' pozz magna' o'vetr, e nun m' fa mal. 
	Vår sære Zulu fra badeøya spilte jo whist og quickstep i min taxi.
	Pa's wijze lynx bezag vroom het fikse aquaduct. 
	Eg kan eta glas utan å skada meg. 
	Vår sære Zulu fra badeøya spilte jo whist og quickstep i min taxi.
	Tsésǫʼ yishą́ągo bííníshghah dóó doo shił neezgai da. 
	Pòdi manjar de veire, me nafrariá pas. 
	ମୁଁ କାଚ ଖାଇପାରେ ଏବଂ ତାହା ମୋର କ୍ଷତି କରିନଥାଏ।. 
	ਮੈਂ ਗਲਾਸ ਖਾ ਸਕਦਾ ਹਾਂ ਅਤੇ ਇਸ ਨਾਲ ਮੈਨੂੰ ਕੋਈ ਤਕਲੀਫ ਨਹੀਂ. 
	Ch'peux mingi du verre, cha m'foé mie n'ma. 
	Pchnąć w tę łódź jeża lub ośm skrzyń fig. 
	Vejam a bruxa da raposa Salta-Pocinhas e o cão feliz que dorme regalado. 
	À noite, vovô Kowalsky vê o ímã cair no pé do pingüim queixoso e vovó põe açúcar no chá de tâmaras do jabuti feliz. 
	Fumegând hipnotic sașiul azvârle mreje în bălți. 
	В чащах юга жил бы цитрус? Да, но фальшивый экземпляр! 
	काचं शक्नोम्यत्तुम् । नोपहिनस्ति माम् ॥ 
	Puotsu mangiari u vitru, nun mi fa mali. 
	මනොපුබ්‌බඞ්‌ගමා ධම්‌මා, මනොසෙට්‌ඨා මනොමයා; මනසා චෙ පදුට්‌ඨෙන, භාසති වා කරොති වා; තතො නං දුක්‌ඛමන්‌වෙති, චක්‌කංව වහතො පදං.
	Starý kôň na hŕbe kníh žuje tíško povädnuté ruže, na stĺpe sa ďateľ učí kvákať novú ódu o živote.
	Šerif bo za vajo spet kuhal domače žgance. 
	Unë mund të ha qelq dhe nuk më gjen gjë. 
	Чешће цeђење мрeжастим џаком побољшава фертилизацију генских хибрида. 
	Flygande bäckasiner söka strax hwila på mjuka tuvor. 
	I kå Glas frässa, ond des macht mr nix! 
	நான் கண்ணாடி சாப்பிடுவேன், அதனால் எனக்கு ஒரு கேடும் வராது. 
	నేను గాజు తినగలను అయినా నాకు యేమీ కాదు. 
	เป็นมนุษย์สุดประเสริฐเลิศคุณค่า - กว่าบรรดาฝูงสัตว์เดรัจฉาน - จงฝ่าฟันพัฒนาวิชาการ อย่าล้างผลาญฤๅเข่นฆ่าบีฑาใคร - ไม่ถือโทษโกรธแช่งซัดฮึดฮัดด่า - หัดอภัยเหมือนกีฬาอัชฌาสัย - ปฏิบัติประพฤติกฎกำหนดใจ - พูดจาให้จ๊ะ ๆ จ๋า ๆ น่าฟังเอยฯ 
	Kaya kong kumain nang bubog at hindi ako masaktan. 
	Pijamalı hasta yağız şoföre çabucak güvendi. 
	Metumi awe tumpan, ɜnyɜ me hwee. 
	Чуєш їх, доцю, га? Кумедна ж ти, прощайся без ґольфів! 
	میں کانچ کھا سکتا ہوں اور مجھے تکلیف نہیں ہوتی ۔ 
	Mi posso magnare el vetro, no'l me fa mae. 
	Con sói nâu nhảy qua con chó lười.
	Dji pou magnî do vêre, çoula m' freut nén må. 
	איך קען עסן גלאָז און עס טוט מיר נישט װײ. 
	Mo lè je̩ dígí, kò ní pa mí lára. 
	我能吞下玻璃而不伤身体。 
	我能吞下玻璃而不傷身體。 
	我能吞下玻璃而不伤身体。 
	我能吞下玻璃而不傷身體。 
	Saya boleh makan kaca dan ia tidak mencederakan saya. 
`
