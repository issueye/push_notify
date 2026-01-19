import { format, formatDistanceToNow } from "date-fns";
import { zhCN } from "date-fns/locale";

export function formatDate(date, pattern = "yyyy-MM-dd HH:mm:ss") {
  if (!date) return "-";
  const d = new Date(date);
  return format(d, pattern, { locale: zhCN });
}

export function formatRelativeTime(date) {
  if (!date) return "-";
  const d = new Date(date);
  return formatDistanceToNow(d, { addSuffix: true, locale: zhCN });
}

export function formatTime(date) {
  return formatDate(date, "HH:mm:ss");
}

export function formatDateOnly(date) {
  return formatDate(date, "YYYY-MM-DD");
}
