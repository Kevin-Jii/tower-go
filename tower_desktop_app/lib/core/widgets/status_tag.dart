import 'package:flutter/material.dart';

class StatusTag extends StatelessWidget {
  final bool active;
  final String enabledText;
  final String disabledText;
  final EdgeInsets padding;
  const StatusTag(
      {super.key,
      required this.active,
      required this.enabledText,
      required this.disabledText,
      this.padding = const EdgeInsets.symmetric(horizontal: 10, vertical: 4)});
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: padding,
      decoration: BoxDecoration(
        color: active ? Colors.green.shade50 : Colors.grey.shade200,
        borderRadius: BorderRadius.circular(4),
        border: Border.all(
            color: active ? Colors.green.shade200 : Colors.grey.shade400),
      ),
      child: Text(
        active ? enabledText : disabledText,
        textAlign: TextAlign.center,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w500,
          color: active ? Colors.green.shade700 : Colors.grey.shade600,
        ),
      ),
    );
  }
}
