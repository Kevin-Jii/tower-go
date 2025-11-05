/// Map 工具：过滤 null，深度兼容简单场景。
library map_utils;

Map<String, dynamic> compact(Map<String, dynamic> src) {
  final result = <String, dynamic>{};
  src.forEach((k, v) {
    if (v != null) result[k] = v;
  });
  return result;
}
