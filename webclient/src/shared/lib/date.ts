export const MONTHS = ["Oca", "Şub", "Mar", "Nis", "May", "Haz", "Tem", "Ağu", "Eyl", "Eki", "Kas", "Ara"];

export function getMonthName(date: Date) {
  return MONTHS[new Date(date).getMonth()];
}

export function formatDate(date: Date, locale: string) {
  return date.toLocaleDateString(locale);
}

export function formatDateString(dateString: string | null, locale: string) {
  if (dateString === null) {
    return "";
  }

  return formatDate(new Date(dateString), locale);
}
