#!/usr/bin/env python3
"""
Simple test script to verify UDC hierarchy logic
"""

def is_numeric(s):
    """Check if a string contains only digits and dots"""
    for c in s:
        if c != '.' and not c.isdigit():
            return False
    return True

def find_numeric_parent(code, prefix, suffix=""):
    """Find the parent for numeric codes with a prefix and optional suffix"""
    if not code:
        return ""
    
    # Handle special ranges like 01/08
    if "/" in code:
        parts = code.split("/")
        if len(parts) > 1:
            # For ranges, the parent is the prefix part
            return prefix + parts[0] + suffix
    
    # Handle numeric codes
    if is_numeric(code):
        if len(code) > 1:
            # Remove last digit to get parent, but preserve dots
            last_char = code[-1]
            if last_char.isdigit():
                parent_code = code[:-1]
                # Remove trailing dot if present
                if parent_code.endswith('.'):
                    parent_code = parent_code[:-1]
                return prefix + parent_code + suffix
    
    return ""

def find_auxiliary_parent(code):
    """Find the parent for auxiliary table codes"""
    if not code:
        return ""
    
    # Handle common auxiliaries of language (=...)
    if code.startswith("=..."):
        if "`" in code:
            # Handle special auxiliary subdivisions like =...`01/`08
            parts = code.split("`")
            if len(parts) > 1:
                return "=..."
        return "TOP"  # =... is a top-level auxiliary
    
    # Handle language codes (=1, =11, =111, etc.)
    if code.startswith("=") and len(code) > 1:
        # Remove the = prefix
        lang_code = code[1:]
        
        # Handle special cases first
        if lang_code in ["00", "030"]:
            return "=..."
        
        # Handle single-digit language codes (top-level)
        if len(lang_code) == 1 and is_numeric(lang_code):
            return "TOP"
        
        # Handle numeric language codes
        if is_numeric(lang_code):
            return find_numeric_parent(lang_code, "=", "")
    
    # Handle place auxiliaries ((1), (1-44), etc.)
    if code.startswith("(") and code.endswith(")"):
        place_code = code[1:-1]
        
        # Handle special cases like (=01)
        if place_code.startswith("="):
            return "(=...)"
        
        # Handle single-digit place codes (top-level)
        if len(place_code) == 1 and is_numeric(place_code):
            return "TOP"
        
        # Handle numeric place codes - for place codes, we need to handle them differently
        # (540) should be a child of (5), not (54)
        if is_numeric(place_code):
            if len(place_code) > 1:
                # For place codes, take the first digit as parent
                first_digit = place_code[0]
                return "(" + first_digit + ")"
    
    # Handle form auxiliaries (-0, -058.6, etc.)
    if code.startswith("-"):
        form_code = code[1:]
        
        # Handle special cases
        if form_code == "0":
            return "TOP"  # -0 is a top-level auxiliary
        
        # Handle numeric form codes
        if is_numeric(form_code):
            return find_numeric_parent(form_code, "-", "")
    
    return ""

def find_main_table_parent(code):
    """Find the parent for main table codes (0-9)"""
    if not code:
        return ""
    
    # Handle special cases
    if code in ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9"]:
        return "TOP"  # Main table divisions are children of TOP
    
    # Handle numeric codes with dots (001, 001.1, etc.)
    if "." in code:
        parts = code.split(".")
        if len(parts) > 1:
            # Remove the last part to get parent
            parent_parts = parts[:-1]
            return ".".join(parent_parts)
    
    # Handle numeric codes without dots (00, 000, 01, etc.)
    if is_numeric(code):
        # Remove last digit to get parent
        if len(code) > 1:
            return code[:-1]
    
    return ""

def find_parent_code(code):
    """Determine the parent UDC code based on the current code"""
    if not code:
        return ""
    
    # Handle auxiliary tables
    if code.startswith("="):
        return find_auxiliary_parent(code)
    if code.startswith("("):
        return find_auxiliary_parent(code)
    if code.startswith("-"):
        return find_auxiliary_parent(code)
    if code.startswith(("+", "/", ":", "[]", "*")) or code == "A/Z":
        return "TOP"  # These are top-level auxiliary signs
    
    # Handle main tables (0-9)
    return find_main_table_parent(code)

def test_hierarchy():
    """Test the hierarchy logic with various UDC codes"""
    test_cases = [
        ("00", "0"),
        ("000", "00"),
        ("001", "00"),
        ("001.1", "001"),
        ("01", "0"),
        ("1", "TOP"),
        ("=1", "TOP"),
        ("=11", "=1"),
        ("=111", "=11"),
        ("(1)", "TOP"),
        ("(540)", "(5)"),
        ("-0", "TOP"),
        ("-058.6", "-058"),
        ("=...", "TOP"),
        ("=...`01", "=..."),
    ]
    
    print("Testing UDC hierarchy logic:")
    print("=" * 40)
    
    all_passed = True
    for code, expected in test_cases:
        result = find_parent_code(code)
        status = "✓" if result == expected else "✗"
        print(f"{status} {code:10} → {result:10} (expected: {expected})")
        if result != expected:
            all_passed = False
    
    print("=" * 40)
    if all_passed:
        print("All tests passed! ✓")
    else:
        print("Some tests failed! ✗")
    
    return all_passed

if __name__ == "__main__":
    test_hierarchy() 