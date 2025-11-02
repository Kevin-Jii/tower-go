/// 统一后端响应格式封装
class ApiResponse<T> {
  final int code;
  final String message;
  final T? data;
  final Map<String, dynamic>? meta; // 兼容分页或附加信息

  ApiResponse({
    required this.code,
    required this.message,
    this.data,
    this.meta,
  });

  bool get success => code == 200;
  bool get hasMeta => meta != null && meta!.isNotEmpty;

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic)? dataParser,
  ) {
    final dynamic rawData = json['data'];
    final parsed = (dataParser != null && rawData != null)
        ? dataParser(rawData)
        : rawData as T?;
    return ApiResponse(
      code: json['code'] ?? 0,
      message: json['message'] ?? '',
      data: parsed,
      meta: json['meta'] is Map<String, dynamic>
          ? (json['meta'] as Map<String, dynamic>)
          : null,
    );
  }
}

/// 分页响应封装
class PageResponse<T> {
  final List<T> list;
  final int total;
  final int page;
  final int pageSize;
  final int? pageCount;
  final bool? hasMore;

  PageResponse({
    required this.list,
    required this.total,
    required this.page,
    required this.pageSize,
    this.pageCount,
    this.hasMore,
  });

  int get totalPages => pageCount ?? (total / pageSize).ceil();
  bool get hasNextPage => hasMore ?? page < totalPages;
  bool get hasPrevPage => page > 1;

  /// 仅从 meta 中解析分页信息（后端约定：meta != null 表示分页接口）
  /// meta 为 null 或缺失关键字段时抛出异常
  factory PageResponse.fromEnvelope(
    Map<String, dynamic> json,
    T Function(Map<String, dynamic>) itemParser,
  ) {
    final meta = json['meta'];
    if (meta is! Map<String, dynamic>) {
      throw ArgumentError('分页响应缺少 meta 或 meta 不是对象');
    }

    final dataList = json['data'];
    if (dataList is! List) {
      throw ArgumentError('期望 data 为 List, 实际为: ${dataList.runtimeType}');
    }

    int requireInt(String key) {
      final v = meta[key];
      if (v is int) return v;
      throw ArgumentError('分页 meta.$key 缺失或类型错误(${v.runtimeType})');
    }

    int? optInt(String key) {
      if (!meta.containsKey(key)) return null;
      final v = meta[key];
      if (v == null) throw ArgumentError('分页 meta.$key 为 null');
      if (v is int) return v;
      throw ArgumentError('分页 meta.$key 类型错误(${v.runtimeType})');
    }

    bool? optBool(String key) {
      if (!meta.containsKey(key)) return null;
      final v = meta[key];
      if (v == null) throw ArgumentError('分页 meta.$key 为 null');
      if (v is bool) return v;
      throw ArgumentError('分页 meta.$key 类型错误(${v.runtimeType})');
    }

    return PageResponse(
      list: dataList.map((e) => itemParser(e as Map<String, dynamic>)).toList(),
      total: requireInt('total'),
      page: requireInt('page'),
      pageSize: requireInt('page_size'),
      pageCount: optInt('page_count'),
      hasMore: optBool('has_more'),
    );
  }
}
