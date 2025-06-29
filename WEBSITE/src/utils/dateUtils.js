/**
 * 格式化日期时间为用户友好的显示格式
 * @param {string} dateString - 日期字符串
 * @returns {string} 格式化后的日期时间字符串
 */
export function formatDate(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = Math.abs(now - date);
  const diffMinutes = Math.floor(diffTime / (1000 * 60));
  const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  // 格式化时间为 HH:MM
  const formatTime = (date) => {
    return date.toLocaleTimeString("zh-CN", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
  };

  // 1分钟内
  if (diffMinutes < 1) {
    return "刚刚";
  }
  // 1小时内
  else if (diffMinutes < 60) {
    return `${diffMinutes}分钟前`;
  }
  // 今天
  else if (diffHours < 24 && date.getDate() === now.getDate()) {
    return `今天 ${formatTime(date)}`;
  }
  // 昨天
  else if (diffDays === 1) {
    return `昨天 ${formatTime(date)}`;
  }
  // 一周内
  else if (diffDays <= 7) {
    const weekdays = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"];
    return `${weekdays[date.getDay()]} ${formatTime(date)}`;
  }
  // 今年内
  else if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString("zh-CN", {
      month: "numeric",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
  }
  // 跨年
  else {
    return date.toLocaleDateString("zh-CN", {
      year: "numeric",
      month: "numeric",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
  }
}

/**
 * 格式化完整的日期时间
 * @param {string} dateString - 日期字符串
 * @returns {string} 完整的日期时间字符串
 */
export function formatFullDateTime(dateString) {
  const date = new Date(dateString);
  return date.toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

/**
 * 格式化为相对时间（多久之前）
 * @param {string} dateString - 日期字符串
 * @returns {string} 相对时间字符串
 */
export function formatRelativeTime(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diffTime = Math.abs(now - date);
  const diffMinutes = Math.floor(diffTime / (1000 * 60));
  const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

  if (diffMinutes < 1) {
    return "刚刚";
  } else if (diffMinutes < 60) {
    return `${diffMinutes}分钟前`;
  } else if (diffHours < 24) {
    return `${diffHours}小时前`;
  } else if (diffDays < 30) {
    return `${diffDays}天前`;
  } else if (diffDays < 365) {
    const months = Math.floor(diffDays / 30);
    return `${months}个月前`;
  } else {
    const years = Math.floor(diffDays / 365);
    return `${years}年前`;
  }
}
