# AYA Web Sitesi

[![Discord](https://img.shields.io/discord/1072074800622739476?color=7289da&logo=discord&logoColor=white)](https://discord.gg/itdepremyardim)
[![GitHub issues](https://img.shields.io/github/issues/eser/aya.is)](https://github.com/eser/aya.is/issues)

(For English please click [here](README.en.md))

Açık kaynak yazılımı ve açık kaynak verilerine ortak ilgi duyan bireylerden oluşan sıkı bağlı ve tutkulu bir topluluğuz.
Toplumun iyileştirilmesi için teknolojiden yararlanma kararlılığımızla birleşerek, Türkiye'de meydana gelen yıkıcı
[Şubat 2023 depremi](https://en.wikipedia.org/wiki/2023_Turkey%E2%80%93Syria_earthquake) sonrasında güçlerimizi
birleştirdik.

Depremin ardından, etkilenenlere destek olma ve aileleri ile sevdiklerini yeniden bir araya getirme sürecine yardım etme
acil ihtiyacını fark ettik. Bu amaç doğrultusunda, bireyler arasındaki arama ve bağlantıyı kolaylaştırmayı hedefleyen
yenilikçi araçlar ve çözümler geliştirmek için yola çıktık.

Yorulmak bilmeyen iş birliği ve kolektif çabalar sayesinde topluluğumuz ilham verici bir büyüklüğe ulaştı ve 24.000
kişiden oluşan bir topluluk haline geldik. Birlikte, aileleri yeniden bir araya getirmeye, umudu geri kazandırmaya ve
depremden etkilenenlere teselli sunmaya yardımcı olmak için teknolojinin gücünden yararlanarak araçlarımızı sürekli
olarak geliştiriyoruz ve genişletiyoruz.

İlerlemeye devam ederken, taahhüdümüzden sapmadan devam ediyoruz. Araçlarımızı geliştirmeye, açık kaynak projelere aktif
olarak katkıda bulunmaya ve bilgi paylaşımı ve iş birliğinin geliştiği destekleyici bir ortamı teşvik etmeye
çalışıyoruz. Açık kaynak topluluğunun bir parçası olmaktan gurur duyuyoruz ve sizlerle olan yolculuğumuza dört gözle
devam etmek için sabırsızlanıyoruz.

## Görevimiz

Açık kaynak çözümlerini kullanarak, bilgi sistemlerini uygulayarak ve mühendislik uygulamalarını hayata geçirerek,
içinde yaşadığımız toplumun iyileştirilmesine aktif olarak katkıda bulunmayı hedefliyoruz. Öncelikli odak noktamız,
sosyal sorumluluk ihtiyaçlarını karşılamak ve bu yollarla çeşitli toplumsal zorluklarla baş etmektir. Uzmanlığımızı ve
kaynaklarımızı kullanarak, topluluğumuz içinde pozitif bir etki yaratmayı ve sürdürülebilir kalkınmayı teşvik etmeyi
amaçlıyoruz.

## Teknoloji

Projelerimizi oluşturmak için kullandığımız teknolojiler şunlardır:

Frontend için:

- [Next.js](https://nextjs.org)
- [Shadcn](https://shadcn/ui)

Önceden kurulu olması gerekenler:

- [Deno](https://deno.land) (önerilen v1.35.2)
- [Node.js](https://nodejs.org) (önerilen v19.0.0)
- [Git](https://git-scm.com/) (önerilen v2.41.0)

## Projeyi Ayağa Kaldırma

GitHub repository'sini klonlayın:

```bash
$ git clone git@github.com:eser/aya.is.git
```

Proje dizinine gidin:

```bash
$ cd aya.is/webclient
```

Gerekli paketleri yükleyin:

```bash
$ deno install --allow-scripts
```

Son olarak, projeyi geliştirme modunda başlatın:

```bash
$ node --run dev
```

## Proje Yönetimi

Şu anda projeye ait bir yönetim panelimiz bulunmamakta. İşlemleri bu nedenle CLI üzerinden yapmaktayız. Sisteminizde
`deno` kurulu ise, aşağıdaki komutlar aracılığı ile CLI'a erişebilirsiniz:

```bash
$ deno task cli
```

CLI üzerinden `backend` nesnesine erişebilirsiniz.

### Örnekler:

Profil getirme:

```js
await backend.getProfile("tr", "eser");
```

## Nasıl Katkıda Bulunabilirsin?

Herkesten katkı bekliyoruz. Başlamak için lütfen [katkıda bulunma kılavuzumuzu](CONTRIBUTING.md) okuyun. Yardım etmek
isterseniz [issue'lara](https://github.com/eser/aya.is/issues) göz atabilirsiniz. Herhangi bir sorunuz varsa
[Discord sunucumuza](https://discord.gg/itdepremyardim) katılmaktan çekinmeyin. Herhangi bir noktada takılırsanız,
GitHub Issues veya Discord üzerinden yardım istemekten çekinmeyin.

## Lisans

Apache 2.0, detaylar için [LICENSE](LICENSE) dosyasını inceleyebilirsiniz.
